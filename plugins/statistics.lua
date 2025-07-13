-- Statistics Plugin for Pomodux
-- Tracks and reports timer usage statistics

-- Register this plugin with Pomodux
pomodux.register_plugin({
    name = "statistics",
    version = "1.0.0",
    description = "Tracks timer usage statistics",
    author = "Pomodux Team"
})

-- Statistics storage
local stats = {
    total_sessions = 0,
    total_work_time = 0,
    total_break_time = 0,
    total_long_break_time = 0,
    completed_sessions = 0,
    interrupted_sessions = 0,
    session_types = {
        work = 0,
        ["break"] = 0,
        long_break = 0
    },
    daily_stats = {}
}

-- Get today's date as a string
function get_today()
    -- This is a simplified version - in a real implementation,
    -- you'd want to use proper date/time functions
    return os.date("%Y-%m-%d")
end

-- Initialize daily stats
function init_daily_stats()
    local today = get_today()
    if not stats.daily_stats[today] then
        stats.daily_stats[today] = {
            work_sessions = 0,
            work_time = 0,
            break_sessions = 0,
            break_time = 0,
            long_break_sessions = 0,
            long_break_time = 0
        }
    end
end

-- Register hooks for timer events
pomodux.register_hook("timer_started", function(event)
    local session_type = event.data.session_type
    local duration = event.data.duration
    
    stats.total_sessions = stats.total_sessions + 1
    stats.session_types[session_type] = stats.session_types[session_type] + 1
    
    init_daily_stats()
    local today = get_today()
    
    if session_type == "work" then
        stats.daily_stats[today].work_sessions = stats.daily_stats[today].work_sessions + 1
    elseif session_type == "break" then
        stats.daily_stats[today].break_sessions = stats.daily_stats[today].break_sessions + 1
    elseif session_type == "long-break" then
        stats.daily_stats[today].long_break_sessions = stats.daily_stats[today].long_break_sessions + 1
    end
    
    pomodux.log(string.format("Statistics: Started %s session (%d minutes)", session_type, duration / 60))
end)

pomodux.register_hook("timer_completed", function(event)
    local session_type = event.data.session_type
    local duration = event.data.duration
    
    stats.completed_sessions = stats.completed_sessions + 1
    
    if session_type == "work" then
        stats.total_work_time = stats.total_work_time + duration
        local today = get_today()
        stats.daily_stats[today].work_time = stats.daily_stats[today].work_time + duration
    elseif session_type == "break" then
        stats.total_break_time = stats.total_break_time + duration
        local today = get_today()
        stats.daily_stats[today].break_time = stats.daily_stats[today].break_time + duration
    elseif session_type == "long-break" then
        stats.total_long_break_time = stats.total_long_break_time + duration
        local today = get_today()
        stats.daily_stats[today].long_break_time = stats.daily_stats[today].long_break_time + duration
    end
    
    pomodux.log(string.format("Statistics: Completed %s session", session_type))
end)

pomodux.register_hook("timer_stopped", function(event)
    stats.interrupted_sessions = stats.interrupted_sessions + 1
    pomodux.log("Statistics: Session interrupted")
end)

-- Function to get statistics
function get_stats()
    return {
        total_sessions = stats.total_sessions,
        completed_sessions = stats.completed_sessions,
        interrupted_sessions = stats.interrupted_sessions,
        total_work_time = stats.total_work_time,
        total_break_time = stats.total_break_time,
        total_long_break_time = stats.total_long_break_time,
        session_types = stats.session_types,
        daily_stats = stats.daily_stats
    }
end

-- Function to get today's statistics
function get_today_stats()
    local today = get_today()
    return stats.daily_stats[today] or {
        work_sessions = 0,
        work_time = 0,
        break_sessions = 0,
        break_time = 0,
        long_break_sessions = 0,
        long_break_time = 0
    }
end

-- Function to format time in minutes
function format_time(seconds)
    return math.floor(seconds / 60)
end

-- Function to print statistics
function print_stats()
    local today_stats = get_today_stats()
    
    pomodux.log("=== POMODUX STATISTICS ===")
    pomodux.log(string.format("Total sessions: %d", stats.total_sessions))
    pomodux.log(string.format("Completed sessions: %d", stats.completed_sessions))
    pomodux.log(string.format("Interrupted sessions: %d", stats.interrupted_sessions))
    pomodux.log(string.format("Total work time: %d minutes", format_time(stats.total_work_time)))
    pomodux.log(string.format("Total break time: %d minutes", format_time(stats.total_break_time)))
    pomodux.log(string.format("Total long break time: %d minutes", format_time(stats.total_long_break_time)))
    pomodux.log("--- Today's Statistics ---")
    pomodux.log(string.format("Work sessions: %d (%d minutes)", today_stats.work_sessions, format_time(today_stats.work_time)))
    pomodux.log(string.format("Break sessions: %d (%d minutes)", today_stats.break_sessions, format_time(today_stats.break_time)))
    pomodux.log(string.format("Long break sessions: %d (%d minutes)", today_stats.long_break_sessions, format_time(today_stats.long_break_time)))
    pomodux.log("=========================")
end

-- Plugin initialization
pomodux.log("Statistics plugin loaded successfully") 