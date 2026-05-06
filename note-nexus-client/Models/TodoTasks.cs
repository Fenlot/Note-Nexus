// Models/TodoTask.cs
using System.Text.Json.Serialization;

namespace note_nexus_client.Models;

public class TodoTask
{
    [JsonPropertyName("id")]
    public long Id { get; set; }
    
    [JsonPropertyName("title")]
    public string Title { get; set; } = string.Empty;
    
    [JsonPropertyName("is_completed")]
    public bool IsCompleted { get; set; }

    [JsonPropertyName("created_at")]
    public DateTime? CreatedAt { get; set; }
}