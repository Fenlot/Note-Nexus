using System.Net.Http.Json;
using System.Threading.Tasks;
using System.Collections.Generic;
using note_nexus_client.Models;

namespace note_nexus_client.Services;

public class NotificationService
{
    private readonly HttpClient _http;

    public NotificationService(HttpClient http)
    {
        _http = http;
    }

    private const string BaseUrl = "v1/notifications";

    public async Task<List<Notification>> GetNotificationsAsync(int limit = 20, int offset = 0)
    {
        try
        {
            return await _http.GetFromJsonAsync<List<Notification>>($"{BaseUrl}?limit={limit}&offset={offset}") ?? new List<Notification>();
        }
        catch
        {
            return new List<Notification>();
        }
    }

    public async Task<int> GetUnreadCountAsync()
    {
        try
        {
            var response = await _http.GetFromJsonAsync<Dictionary<string, int>>($"{BaseUrl}/unread");
            return response != null && response.ContainsKey("unread_count") ? response["unread_count"] : 0;
        }
        catch
        {
            return 0;
        }
    }

    public async Task MarkAsReadAsync(long notificationId)
    {
        try
        {
            var request = new { is_read = true };
            await _http.PatchAsJsonAsync($"{BaseUrl}/{notificationId}", request);
        }
        catch
        {
            throw;
        }
    }

    public async Task MarkAllAsReadAsync()
    {
        try
        {
            await _http.PostAsJsonAsync($"{BaseUrl}/mark-all-read", new { });
        }
        catch
        {
            throw;
        }
    }

    public async Task DeleteNotificationAsync(long notificationId)
    {
        try
        {
            await _http.DeleteAsync($"{BaseUrl}/{notificationId}");
        }
        catch
        {
            throw;
        }
    }
}
