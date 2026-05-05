// Models/TodoTask.cs
namespace note_nexus_client.Models;

public class TodoTask
{
    // These names MUST match the JSON 'key' from your Go backend exactly.
    public long Id { get; set; }
    
    public string Title { get; set; } = string.Empty;
    
    public bool IsCompleted { get; set; }
}