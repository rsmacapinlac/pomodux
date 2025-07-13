-- Mako Notification Plugin for Pomodux
-- Provides system notifications using mako (Wayland notification daemon)

pomodux.register_plugin({
    name = "mako_notification",
    version = "1.0.0",
    description = "Provides system notifications using mako",
    author = "Pomodux Team"
})

local function send_notification(title, message)
    -- Use notify-send which will be handled by mako on Wayland
    local notify_cmd = string.format("notify-send '%s' '%s'", title, message)
    os.execute(notify_cmd)
end

pomodux.register_hook("timer_completed", function(event)
    print("LUA DEBUG: timer_completed hook triggered")
    pomodux.log("DEBUG: timer_completed hook triggered")
    local session_type = event.data.session_type
    print("LUA DEBUG: session_type is " .. tostring(session_type))
    pomodux.log("DEBUG: session_type is " .. tostring(session_type))
    local title, message
    if session_type == "work" then
        title = "Work Session Complete"
        message = "Your work session has finished. Time for a break!"
    elseif session_type == "break" then
        title = "Break Complete"
        message = "Break time is over. Ready to work?"
    elseif session_type == "long-break" then
        title = "Long Break Complete"
        message = "Long break finished. Time to get back to work!"
    else
        title = "Timer Complete"
        message = "Your timer session has finished."
    end
    print("LUA DEBUG: Sending completion notification: " .. title .. " - " .. message)
    pomodux.log("DEBUG: Sending completion notification: " .. title .. " - " .. message)
    send_notification(title, message)
    print("LUA DEBUG: send_notification called")
end)

pomodux.register_hook("timer_started", function(event)
    local session_type = event.data.session_type
    local duration = event.data.duration
    pomodux.log("DEBUG: duration value is " .. tostring(duration) .. " (type: " .. type(duration) .. ")")
    local title = "Timer Started"
    local message
    if duration < 60 then
        message = string.format("%s session started for %d seconds", session_type, duration)
    else
        local minutes = math.floor(duration / 60)
        message = string.format("%s session started for %d minutes", session_type, minutes)
    end
    send_notification(title, message)
end)

pomodux.register_hook("timer_paused", function(event)
    local title = "Timer Paused"
    local message = "Your timer session has been paused."
    send_notification(title, message)
end)

pomodux.register_hook("timer_resumed", function(event)
    local title = "Timer Resumed"
    local message = "Your timer session has been resumed."
    send_notification(title, message)
end)

pomodux.register_hook("timer_stopped", function(event)
    local title = "Timer Stopped"
    local message = "Your timer session has been stopped."
    send_notification(title, message)
end) 