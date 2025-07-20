package plugin

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/rsmacapinlac/pomodux/internal/logger"

	lua "github.com/yuin/gopher-lua"
)

// EventType represents the type of timer event
type EventType string

const (
	EventTimerStarted   EventType = "timer_started"
	EventTimerPaused    EventType = "timer_paused"
	EventTimerResumed   EventType = "timer_resumed"
	EventTimerCompleted EventType = "timer_completed"
	EventTimerStopped   EventType = "timer_stopped"
)

// Event represents a timer event
type Event struct {
	Type      EventType
	Timestamp time.Time
	Data      map[string]interface{}
}

// Plugin represents a loaded plugin
type Plugin struct {
	Name        string
	Version     string
	Description string
	Author      string
	LState      *lua.LState
	Hooks       map[EventType][]lua.LValue
	Enabled     bool
	mu          sync.RWMutex
}

// PluginManager manages the plugin system
type PluginManager struct {
	plugins    map[string]*Plugin
	events     chan Event
	mu         sync.RWMutex
	done       chan struct{}
	pluginsDir string
	api        *PluginAPI
}

// PluginAPI provides the interface for plugins to register themselves
type PluginAPI struct {
	manager *PluginManager
}

// NewPluginManager creates a new plugin manager
func NewPluginManager(pluginsDir string) *PluginManager {
	pm := &PluginManager{
		plugins:    make(map[string]*Plugin),
		events:     make(chan Event, 100),
		done:       make(chan struct{}),
		pluginsDir: pluginsDir,
	}

	pm.api = &PluginAPI{manager: pm}

	// Start event processing goroutine
	go pm.processEvents()

	return pm
}

// LoadPlugins loads all plugins from the plugins directory
func (pm *PluginManager) LoadPlugins() error {
	// Create plugins directory if it doesn't exist
	if err := os.MkdirAll(pm.pluginsDir, 0750); err != nil {
		return fmt.Errorf("failed to create plugins directory: %w", err)
	}

	// Walk through plugins directory
	return filepath.WalkDir(pm.pluginsDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip directories and non-lua files
		if d.IsDir() || filepath.Ext(path) != ".lua" {
			return nil
		}

		// Load the plugin
		if err := pm.LoadPluginFromFile(path); err != nil {
			logger.Warn("Failed to load plugin", map[string]interface{}{"path": path, "error": err.Error()})
			return nil // Continue loading other plugins
		}

		return nil
	})
}

// LoadPluginFromFile loads a plugin from a Lua file
func (pm *PluginManager) LoadPluginFromFile(filePath string) error {
	// Validate file path for security
	if err := validateFilePath(filePath); err != nil {
		return fmt.Errorf("invalid file path: %w", err)
	}

	// Read the plugin file
	content, err := os.ReadFile(filePath) // #nosec G304 -- filePath is validated by validateFilePath
	if err != nil {
		return fmt.Errorf("failed to read plugin file %s: %w", filePath, err)
	}

	// Extract plugin name from filename
	pluginName := filepath.Base(filePath)
	pluginName = pluginName[:len(pluginName)-4] // Remove .lua extension

	return pm.LoadPlugin(pluginName, string(content))
}

// LoadPlugin loads a plugin from Lua code
func (pm *PluginManager) LoadPlugin(name, code string) error {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	// Check if plugin already exists
	if _, exists := pm.plugins[name]; exists {
		return fmt.Errorf("plugin %s already loaded", name)
	}

	// Create new Lua state for the plugin
	L := lua.NewState()

	// Register the plugin API
	pm.registerPluginAPI(L)

	// Load and run the plugin code
	if err := L.DoString(code); err != nil {
		L.Close()
		return fmt.Errorf("failed to load plugin %s: %w", name, err)
	}

	// Check if plugin registered itself
	pluginTable := L.GetGlobal("plugin")
	if pluginTable.Type() != lua.LTTable {
		L.Close()
		return fmt.Errorf("plugin %s must register itself using pomodux.register_plugin()", name)
	}

	// Extract plugin info
	pluginInfo := pluginTable.(*lua.LTable)
	version := L.GetField(pluginInfo, "version")
	description := L.GetField(pluginInfo, "description")
	author := L.GetField(pluginInfo, "author")

	// Create plugin instance
	plugin := &Plugin{
		Name:        name,
		Version:     version.String(),
		Description: description.String(),
		Author:      author.String(),
		LState:      L,
		Hooks:       make(map[EventType][]lua.LValue),
		Enabled:     true,
	}

	// Merge any pending hooks from Lua state
	if hooksTable := L.GetGlobal("__pomodux_pending_hooks"); hooksTable.Type() == lua.LTTable {
		hooks := hooksTable.(*lua.LTable)
		hooks.ForEach(func(key lua.LValue, value lua.LValue) {
			eventType := EventType(key.String())
			if value.Type() == lua.LTFunction {
				plugin.Hooks[eventType] = append(plugin.Hooks[eventType], value)
			}
		})
		L.SetGlobal("__pomodux_pending_hooks", lua.LNil)
	}

	// Store plugin
	pm.plugins[name] = plugin

	logger.Info("Loaded plugin", map[string]interface{}{"name": name, "version": version, "author": author})
	return nil
}

// registerPluginAPI registers the plugin API functions in the Lua state
func (pm *PluginManager) registerPluginAPI(L *lua.LState) {
	// Create the pomodux table
	pomoduxTable := L.CreateTable(0, 2)
	L.SetGlobal("pomodux", pomoduxTable)

	// Register plugin registration function
	registerFn := L.NewFunction(func(L *lua.LState) int {
		pluginInfo := L.CheckTable(1)
		pm.api.registerPlugin(L, pluginInfo)
		return 0
	})
	pomoduxTable.RawSetString("register_plugin", registerFn)

	// Register hook registration function
	hookFn := L.NewFunction(func(L *lua.LState) int {
		eventType := L.CheckString(1)
		callback := L.CheckFunction(2)
		pm.api.registerHook(L, eventType, callback)
		return 0
	})
	pomoduxTable.RawSetString("register_hook", hookFn)

	// Register utility functions
	pm.registerUtilityFunctions(L, pomoduxTable)
}

// registerUtilityFunctions registers utility functions for plugins
func (pm *PluginManager) registerUtilityFunctions(L *lua.LState, pomoduxTable *lua.LTable) {
	// Log function
	logFn := L.NewFunction(func(L *lua.LState) int {
		message := L.CheckString(1)
		logger.Debug("[PLUGIN] " + message)
		return 0
	})
	pomoduxTable.RawSetString("log", logFn)

	// Get config function
	getConfigFn := L.NewFunction(func(L *lua.LState) int {
		_ = L.CheckString(1) // key parameter (not used yet)
		defaultValue := L.OptString(2, "")
		// TODO: Implement config retrieval
		L.Push(lua.LString(defaultValue))
		return 1
	})
	pomoduxTable.RawSetString("get_config", getConfigFn)
}

// registerPlugin registers a plugin with the manager
func (api *PluginAPI) registerPlugin(L *lua.LState, pluginInfo *lua.LTable) {
	// Store the plugin info globally so the manager can access it
	L.SetGlobal("plugin", pluginInfo)

	// If there are any hooks registered before plugin registration, move them to the plugin.Hooks map
	if hooksTable := L.GetGlobal("__pomodux_pending_hooks"); hooksTable.Type() == lua.LTTable {
		// Find the plugin in the manager
		for _, plugin := range api.manager.plugins {
			if plugin.LState == L {
				plugin.mu.Lock()
				hooks := hooksTable.(*lua.LTable)
				hooks.ForEach(func(key lua.LValue, value lua.LValue) {
					eventType := EventType(key.String())
					if value.Type() == lua.LTFunction {
						plugin.Hooks[eventType] = append(plugin.Hooks[eventType], value)
					}
				})
				plugin.mu.Unlock()
			}
		}
		L.SetGlobal("__pomodux_pending_hooks", lua.LNil)
	}
}

// registerHook registers a hook for a specific event type
func (api *PluginAPI) registerHook(L *lua.LState, eventType string, callback *lua.LFunction) {
	// Find the current plugin
	pluginName := ""
	for name, plugin := range api.manager.plugins {
		if plugin.LState == L {
			pluginName = name
			break
		}
	}

	if pluginName == "" {
		// Plugin not yet registered, store in a temporary table in Lua state
		hooksTable := L.GetGlobal("__pomodux_pending_hooks")
		if hooksTable.Type() != lua.LTTable {
			hooksTable = L.CreateTable(0, 5)
			L.SetGlobal("__pomodux_pending_hooks", hooksTable)
		}
		hooksTable.(*lua.LTable).RawSetString(eventType, callback)
		return
	}

	// Register the hook
	plugin := api.manager.plugins[pluginName]
	eventTypeEnum := EventType(eventType)
	plugin.mu.Lock()
	plugin.Hooks[eventTypeEnum] = append(plugin.Hooks[eventTypeEnum], callback)
	plugin.mu.Unlock()

	logger.Info("Registered hook", map[string]interface{}{"event_type": eventType, "plugin": pluginName})
}

// EmitEvent emits an event to all registered plugins
func (pm *PluginManager) EmitEvent(event Event) {
	logger.Debug("PLUGIN: EmitEvent called", map[string]interface{}{"event": event.Type})
	select {
	case pm.events <- event:
		logger.Debug("PLUGIN: Event queued for processing", map[string]interface{}{"event": event.Type})
	default:
		logger.Warn("PLUGIN: Event channel full, dropping event", map[string]interface{}{"event": event.Type})
	}
}

// processEvents processes events and calls plugin hooks
func (pm *PluginManager) processEvents() {
	for {
		select {
		case event := <-pm.events:
			pm.callPluginHooks(event)
		case <-pm.done:
			return
		}
	}
}

// callPluginHooks calls all registered hooks for an event
func (pm *PluginManager) callPluginHooks(event Event) {
	logger.Debug("PLUGIN: callPluginHooks", map[string]interface{}{"event": event.Type})
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	for _, plugin := range pm.plugins {
		if !plugin.Enabled {
			continue
		}

		plugin.mu.RLock()
		hooks := plugin.Hooks[event.Type]
		plugin.mu.RUnlock()

		logger.Debug("PLUGIN: Plugin has hooks", map[string]interface{}{"plugin": plugin.Name, "hook_count": len(hooks), "event": event.Type})
		for _, hook := range hooks {
			go func(p *Plugin, h lua.LValue) {
				logger.Debug("PLUGIN: Calling hook", map[string]interface{}{"plugin": p.Name, "event": event.Type})
				if err := pm.callHook(p, h, event); err != nil {
					logger.Error("PLUGIN: Error calling hook", err, map[string]interface{}{"plugin": p.Name})
				}
			}(plugin, hook)
		}
	}
}

// callHook calls a single plugin hook
func (pm *PluginManager) callHook(plugin *Plugin, hook lua.LValue, event Event) error {
	L := plugin.LState

	// Create event table for Lua
	eventTable := L.CreateTable(0, 3)
	eventTable.RawSetString("type", lua.LString(event.Type))
	eventTable.RawSetString("timestamp", lua.LNumber(event.Timestamp.Unix()))

	// Create data table
	dataTable := L.CreateTable(0, len(event.Data))
	for k, v := range event.Data {
		switch val := v.(type) {
		case string:
			dataTable.RawSetString(k, lua.LString(val))
		case int:
			dataTable.RawSetString(k, lua.LNumber(val))
		case float64:
			dataTable.RawSetString(k, lua.LNumber(val))
		case bool:
			dataTable.RawSetString(k, lua.LBool(val))
		default:
			dataTable.RawSetString(k, lua.LString(fmt.Sprintf("%v", val)))
		}
	}
	eventTable.RawSetString("data", dataTable)

	// Call the hook function
	logger.Debug("PLUGIN: About to call Lua hook", map[string]interface{}{"plugin": plugin.Name, "event": event.Type})
	if err := L.CallByParam(lua.P{
		Fn:      hook.(*lua.LFunction),
		NRet:    0,
		Protect: true,
	}, eventTable); err != nil {
		logger.Error("PLUGIN: Error calling Lua hook", err, map[string]interface{}{"plugin": plugin.Name, "event": event.Type})
		return fmt.Errorf("failed to call hook: %w", err)
	}
	logger.Debug("PLUGIN: Successfully called Lua hook", map[string]interface{}{"plugin": plugin.Name, "event": event.Type})

	return nil
}

// GetPlugin returns a plugin by name
func (pm *PluginManager) GetPlugin(name string) (*Plugin, bool) {
	pm.mu.RLock()
	defer pm.mu.RUnlock()
	plugin, exists := pm.plugins[name]
	return plugin, exists
}

// ListPlugins returns all loaded plugins
func (pm *PluginManager) ListPlugins() []*Plugin {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	plugins := make([]*Plugin, 0, len(pm.plugins))
	for _, plugin := range pm.plugins {
		plugins = append(plugins, plugin)
	}
	return plugins
}

// EnablePlugin enables or disables a plugin
func (pm *PluginManager) EnablePlugin(name string, enabled bool) error {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	plugin, exists := pm.plugins[name]
	if !exists {
		return fmt.Errorf("plugin %s not found", name)
	}

	plugin.Enabled = enabled
	return nil
}

// UnloadPlugin unloads a plugin
func (pm *PluginManager) UnloadPlugin(name string) error {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	plugin, exists := pm.plugins[name]
	if !exists {
		return fmt.Errorf("plugin %s not found", name)
	}

	// Close the Lua state
	plugin.LState.Close()

	// Remove from plugins map
	delete(pm.plugins, name)

	return nil
}

// Shutdown shuts down the plugin manager
func (pm *PluginManager) Shutdown() {
	close(pm.done)

	// Unload all plugins
	pm.mu.Lock()
	defer pm.mu.Unlock()

	for name, plugin := range pm.plugins {
		plugin.LState.Close()
		delete(pm.plugins, name)
	}
}

// validateFilePath validates a file path for security
func validateFilePath(filePath string) error {
	// Check for path traversal attempts
	if strings.Contains(filePath, "..") {
		return fmt.Errorf("path traversal not allowed")
	}

	// Check for dangerous characters
	dangerousChars := []string{"|", "&", ";", "`", "$", "(", ")", "<", ">", "*", "?"}
	for _, char := range dangerousChars {
		if strings.Contains(filePath, char) {
			return fmt.Errorf("dangerous character '%s' not allowed in file path", char)
		}
	}

	// For absolute paths, check if they're in a safe location
	if filepath.IsAbs(filePath) {
		// Allow paths in /tmp for testing
		if strings.HasPrefix(filePath, "/tmp/") {
			return nil
		}

		// Allow paths in the user's home directory
		homeDir, err := os.UserHomeDir()
		if err == nil {
			if strings.HasPrefix(filePath, homeDir) {
				return nil
			}
		}

		// For now, be restrictive with other absolute paths
		return fmt.Errorf("absolute paths not allowed for security")
	}

	return nil
}
