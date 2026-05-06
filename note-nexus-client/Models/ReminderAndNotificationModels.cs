namespace note_nexus_client.Models;

public class Reminder
{
    public long Id { get; set; }
    public long WorkspaceId { get; set; }
    public long UserId { get; set; }
    public long? TargetId { get; set; }
    public string TargetType { get; set; } // 'task', 'note', 'custom'
    public string ReminderType { get; set; } // 'task_due', 'note_review', 'custom'
    public string Title { get; set; }
    public string Description { get; set; }
    public DateTime DueDate { get; set; }
    public string ScheduleType { get; set; } // 'once', 'daily', 'weekly', 'monthly'
    public bool IsActive { get; set; }
    public DateTime CreatedAt { get; set; }
    public DateTime UpdatedAt { get; set; }
}

public class Notification
{
    public long Id { get; set; }
    public long UserId { get; set; }
    public long? ReminderId { get; set; }
    public string Title { get; set; }
    public string Message { get; set; }
    public string NotificationType { get; set; } // 'reminder', 'task_update', 'mention'
    public bool IsRead { get; set; }
    public string Channels { get; set; } // JSON
    public DateTime CreatedAt { get; set; }
    public DateTime? ReadAt { get; set; }
}
