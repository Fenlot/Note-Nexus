using System.Net.Http.Json;
using System.Threading.Tasks;
using System.Collections.Generic;
using note_nexus_client.Models;

namespace note_nexus_client.Services;

public class ReminderService
{
    private readonly HttpClient _http;
    private readonly AppState _appState;

    public ReminderService(HttpClient http, AppState appState)
    {
        _http = http;
        _appState = appState;
    }

    private string BaseUrl => $"v1/workspaces/{_appState.ActiveWorkspaceId}/reminders";

    public async Task<List<Reminder>> GetRemindersAsync()
    {
        if (_appState.ActiveWorkspaceId == 0) return new List<Reminder>();
        try
        {
            return await _http.GetFromJsonAsync<List<Reminder>>(BaseUrl) ?? new List<Reminder>();
        }
        catch
        {
            return new List<Reminder>();
        }
    }

    public async Task<Reminder> GetReminderAsync(long id)
    {
        if (_appState.ActiveWorkspaceId == 0) return null;
        try
        {
            return await _http.GetFromJsonAsync<Reminder>($"{BaseUrl}/{id}");
        }
        catch
        {
            return null;
        }
    }

    public async Task CreateReminderAsync(Reminder reminder)
    {
        if (_appState.ActiveWorkspaceId == 0) return;
        
        var request = new
        {
            target_id = reminder.TargetId,
            target_type = reminder.TargetType,
            reminder_type = reminder.ReminderType,
            title = reminder.Title,
            description = reminder.Description,
            due_date = reminder.DueDate.ToString("O"), // RFC3339 format
            schedule_type = reminder.ScheduleType
        };

        try
        {
            await _http.PostAsJsonAsync(BaseUrl, request);
        }
        catch
        {
            throw;
        }
    }

    public async Task UpdateReminderAsync(long id, Reminder reminder)
    {
        if (_appState.ActiveWorkspaceId == 0) return;

        var request = new
        {
            title = reminder.Title,
            description = reminder.Description,
            due_date = reminder.DueDate.ToString("O"),
            schedule_type = reminder.ScheduleType,
            is_active = reminder.IsActive
        };

        try
        {
            await _http.PutAsJsonAsync($"{BaseUrl}/{id}", request);
        }
        catch
        {
            throw;
        }
    }

    public async Task DeleteReminderAsync(long id)
    {
        if (_appState.ActiveWorkspaceId == 0) return;

        try
        {
            await _http.DeleteAsync($"{BaseUrl}/{id}");
        }
        catch
        {
            throw;
        }
    }
}
