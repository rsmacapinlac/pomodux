-- Debug Events Plugin for Pomodux
-- Prints messages for all timer events to help debug the plugin system

pomodux.register_plugin({
    name = "debug_events",
    version = "1.0.0",
    description = "Debug plugin that prints messages for all timer events",
    author = "Pomodux Team"
})

pomodux.register_hook("timer_started", function(event)
    local session_type = event.data.session_type
    local duration = event.data.duration
    print("🔵 DEBUG: timer_started - " .. session_type .. " session for " .. duration .. " seconds")
end)

pomodux.register_hook("timer_completed", function(event)
    local session_type = event.data.session_type
    local duration = event.data.duration
    print("🟢 DEBUG: timer_completed - " .. session_type .. " session completed after " .. duration .. " seconds")
end)

pomodux.register_hook("timer_stopped", function(event)
    local session_type = event.data.session_type
    local duration = event.data.duration
    print("🔴 DEBUG: timer_stopped - " .. session_type .. " session stopped after " .. duration .. " seconds")
end)

pomodux.register_hook("timer_paused", function(event)
    local session_type = event.data.session_type
    local elapsed = event.data.elapsed
    print("⏸️  DEBUG: timer_paused - " .. session_type .. " session paused at " .. elapsed .. " seconds")
end)

pomodux.register_hook("timer_resumed", function(event)
    local session_type = event.data.session_type
    local elapsed = event.data.elapsed
    print("▶️  DEBUG: timer_resumed - " .. session_type .. " session resumed at " .. elapsed .. " seconds")
end) 