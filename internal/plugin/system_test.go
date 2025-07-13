package plugin

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/rsmacapinlac/pomodux/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPluginSystem_FileLoading(t *testing.T) {
	pluginsDir := t.TempDir()
	pm := NewPluginManager(pluginsDir)
	defer pm.Shutdown()

	// Create a test plugin file
	pluginFile := filepath.Join(pluginsDir, "test_plugin.lua")
	pluginCode := `
pomodux.register_plugin({
    name = "test_plugin",
    version = "1.0.0",
    description = "Test plugin for file loading",
    author = "Test Author"
})

pomodux.register_hook("timer_started", function(event)
    pomodux.log("Timer started: " .. event.data.session_type)
end)

pomodux.register_hook("timer_completed", function(event)
    pomodux.log("Timer completed: " .. event.data.session_type)
end)
`

	err := os.WriteFile(pluginFile, []byte(pluginCode), 0644)
	require.NoError(t, err)

	// Load plugins from directory
	err = pm.LoadPlugins()
	require.NoError(t, err)

	// Verify plugin was loaded
	plugin, exists := pm.GetPlugin("test_plugin")
	assert.True(t, exists)
	assert.Equal(t, "test_plugin", plugin.Name)
	assert.Equal(t, "1.0.0", plugin.Version)
	assert.Equal(t, "Test plugin for file loading", plugin.Description)
	assert.Equal(t, "Test Author", plugin.Author)
	assert.True(t, plugin.Enabled)

	// Verify hooks were registered
	assert.NotEmpty(t, plugin.Hooks[EventTimerStarted])
	assert.NotEmpty(t, plugin.Hooks[EventTimerCompleted])
}

func TestPluginSystem_EventEmission(t *testing.T) {
	pluginsDir := t.TempDir()
	pm := NewPluginManager(pluginsDir)
	defer pm.Shutdown()

	// Create a test plugin that logs events
	pluginFile := filepath.Join(pluginsDir, "event_logger.lua")
	pluginCode := `
local event_count = 0

pomodux.register_plugin({
    name = "event_logger",
    version = "1.0.0",
    description = "Logs all events",
    author = "Test Author"
})

pomodux.register_hook("timer_started", function(event)
    event_count = event_count + 1
    pomodux.log("Event received: " .. event.type .. " (count: " .. event_count .. ")")
end)

pomodux.register_hook("timer_completed", function(event)
    event_count = event_count + 1
    pomodux.log("Event received: " .. event.type .. " (count: " .. event_count .. ")")
end)
`

	err := os.WriteFile(pluginFile, []byte(pluginCode), 0644)
	require.NoError(t, err)

	// Load plugins
	err = pm.LoadPlugins()
	require.NoError(t, err)

	// Emit events
	events := []Event{
		{
			Type:      EventTimerStarted,
			Timestamp: time.Now(),
			Data: map[string]interface{}{
				"session_type": "work",
				"duration":     1500,
			},
		},
		{
			Type:      EventTimerCompleted,
			Timestamp: time.Now(),
			Data: map[string]interface{}{
				"session_type": "work",
				"duration":     1500,
			},
		},
	}

	for _, event := range events {
		pm.EmitEvent(event)
		time.Sleep(50 * time.Millisecond) // Give time for event processing
	}
}

func TestPluginSystem_PluginManagement(t *testing.T) {
	pluginsDir := t.TempDir()
	pm := NewPluginManager(pluginsDir)
	defer pm.Shutdown()

	// Create multiple test plugins
	plugins := []struct {
		name string
		code string
	}{
		{
			name: "plugin1.lua",
			code: `
pomodux.register_plugin({
    name = "plugin1",
    version = "1.0.0",
    description = "First test plugin",
    author = "Test Author"
})
`,
		},
		{
			name: "plugin2.lua",
			code: `
pomodux.register_plugin({
    name = "plugin2",
    version = "1.0.0",
    description = "Second test plugin",
    author = "Test Author"
})
`,
		},
	}

	// Create plugin files
	for _, p := range plugins {
		pluginFile := filepath.Join(pluginsDir, p.name)
		err := os.WriteFile(pluginFile, []byte(p.code), 0644)
		require.NoError(t, err)
	}

	// Load plugins
	err := pm.LoadPlugins()
	require.NoError(t, err)

	// Test plugin listing
	loadedPlugins := pm.ListPlugins()
	assert.Len(t, loadedPlugins, 2)

	// Test plugin enable/disable
	err = pm.EnablePlugin("plugin1", false)
	require.NoError(t, err)

	plugin1, exists := pm.GetPlugin("plugin1")
	assert.True(t, exists)
	assert.False(t, plugin1.Enabled)

	err = pm.EnablePlugin("plugin1", true)
	require.NoError(t, err)

	plugin1, _ = pm.GetPlugin("plugin1")
	assert.True(t, plugin1.Enabled)

	// Test plugin unloading
	err = pm.UnloadPlugin("plugin2")
	require.NoError(t, err)

	_, exists = pm.GetPlugin("plugin2")
	assert.False(t, exists)

	loadedPlugins = pm.ListPlugins()
	assert.Len(t, loadedPlugins, 1)
}

func TestPluginSystem_InvalidPlugins(t *testing.T) {
	pluginsDir := t.TempDir()
	pm := NewPluginManager(pluginsDir)
	defer pm.Shutdown()

	// Test plugin without registration
	pluginFile := filepath.Join(pluginsDir, "invalid_plugin.lua")
	pluginCode := `
-- This plugin doesn't register itself
print("Hello world")
`

	err := os.WriteFile(pluginFile, []byte(pluginCode), 0644)
	require.NoError(t, err)

	// Should fail to load
	err = pm.LoadPlugins()
	require.NoError(t, err) // LoadPlugins continues on individual plugin errors

	// Plugin should not be loaded
	_, exists := pm.GetPlugin("invalid_plugin")
	assert.False(t, exists)
}

func TestPluginSystem_ConfigBasedDirectory(t *testing.T) {
	// Test that plugin manager uses config-based directory
	cfg := config.DefaultConfig()
	pluginsDir := cfg.Plugins.Directory

	pm := NewPluginManager(pluginsDir)
	defer pm.Shutdown()

	// Verify the plugin manager uses the config directory
	assert.Equal(t, pluginsDir, pm.pluginsDir)
}

func TestPluginSystem_ConcurrentEventEmission(t *testing.T) {
	pluginsDir := t.TempDir()
	pm := NewPluginManager(pluginsDir)
	defer pm.Shutdown()

	// Create a test plugin
	pluginFile := filepath.Join(pluginsDir, "concurrent_test.lua")
	pluginCode := `
pomodux.register_plugin({
    name = "concurrent_test",
    version = "1.0.0",
    description = "Test concurrent event handling",
    author = "Test Author"
})

pomodux.register_hook("timer_started", function(event)
    pomodux.log("Concurrent event: " .. event.type)
end)
`

	err := os.WriteFile(pluginFile, []byte(pluginCode), 0644)
	require.NoError(t, err)

	// Load plugins
	err = pm.LoadPlugins()
	require.NoError(t, err)

	// Emit multiple events concurrently
	for i := 0; i < 10; i++ {
		go func(id int) {
			event := Event{
				Type:      EventTimerStarted,
				Timestamp: time.Now(),
				Data: map[string]interface{}{
					"session_type": "work",
					"duration":     1500,
					"id":           id,
				},
			}
			pm.EmitEvent(event)
		}(i)
	}

	// Give time for all events to be processed
	time.Sleep(200 * time.Millisecond)
}

func TestPluginSystem_Shutdown(t *testing.T) {
	pluginsDir := t.TempDir()
	pm := NewPluginManager(pluginsDir)

	// Create and load a plugin
	pluginFile := filepath.Join(pluginsDir, "shutdown_test.lua")
	pluginCode := `
pomodux.register_plugin({
    name = "shutdown_test",
    version = "1.0.0",
    description = "Test plugin shutdown",
    author = "Test Author"
})
`

	err := os.WriteFile(pluginFile, []byte(pluginCode), 0644)
	require.NoError(t, err)

	err = pm.LoadPlugins()
	require.NoError(t, err)

	// Verify plugin is loaded
	_, exists := pm.GetPlugin("shutdown_test")
	assert.True(t, exists)

	// Shutdown plugin manager
	pm.Shutdown()

	// Verify plugin is no longer accessible
	_, exists = pm.GetPlugin("shutdown_test")
	assert.False(t, exists)
}
