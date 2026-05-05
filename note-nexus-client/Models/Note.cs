using System.Text.Json.Serialization; // Need this for the attribute

namespace note_nexus_client.Models;

public class Note
{
    [JsonPropertyName("id")]
    public long Id { get; set; }

    [JsonPropertyName("title")]
    public string Title { get; set; } = string.Empty;

    [JsonPropertyName("content")]
    public string Content { get; set; } = string.Empty;
}