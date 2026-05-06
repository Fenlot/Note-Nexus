using System.Text.Json.Serialization;

namespace note_nexus_client.Models;

public class Reminder
{
    [JsonPropertyName("id")]
    public long Id { get; set; }

    [JsonPropertyName("workspace_id")]
    public long WorkspaceId { get; set; }

    [JsonPropertyName("user_id")]
    public long UserId { get; set; }

    [JsonPropertyName("target_id")]
    public long? TargetId { get; set; }

    [JsonPropertyName("target_type")]
    public string TargetType { get; set; } // 'task', 'note', 'custom'

    [JsonPropertyName("reminder_type")]
    public string ReminderType { get; set; } // 'task_due', 'note_review', 'custom'

    [JsonPropertyName("title")]
    public string Title { get; set; }

    [JsonPropertyName("description")]
    public string Description { get; set; }

    [JsonPropertyName("due_date")]
    public DateTime DueDate { get; set; }

    [JsonPropertyName("schedule_type")]
    public string ScheduleType { get; set; } // 'once', 'daily', 'weekly', 'monthly'

    [JsonPropertyName("is_active")]
    public bool IsActive { get; set; }

    [JsonPropertyName("created_at")]
    public DateTime CreatedAt { get; set; }

    [JsonPropertyName("updated_at")]
    public DateTime UpdatedAt { get; set; }
}

public class Notification
{
    [JsonPropertyName("id")]
    public long Id { get; set; }

    [JsonPropertyName("user_id")]
    public long UserId { get; set; }

    [JsonPropertyName("reminder_id")]
    public long? ReminderId { get; set; }

    [JsonPropertyName("title")]
    public string Title { get; set; }

    [JsonPropertyName("message")]
    public string Message { get; set; }

    [JsonPropertyName("notification_type")]
    public string NotificationType { get; set; } // 'reminder', 'task_update', 'mention'

    [JsonPropertyName("is_read")]
    public bool IsRead { get; set; }

    [JsonPropertyName("channels")]
    public string Channels { get; set; } // JSON

    [JsonPropertyName("created_at")]
    public DateTime CreatedAt { get; set; }

    [JsonPropertyName("read_at")]
    public DateTime? ReadAt { get; set; }
}
