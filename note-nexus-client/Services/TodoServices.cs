using System.Net.Http.Json;
using System.Threading.Tasks;
using System.Collections.Generic;
using note_nexus_client.Models;

namespace note_nexus_client.Services;

public class TodoService
{
    private readonly HttpClient _http;
    private readonly AppState _appState;

    public TodoService(HttpClient http, AppState appState)
    {
        _http = http;
        _appState = appState;
    }

    private string BaseUrl => $"v1/workspaces/{_appState.ActiveWorkspaceId}/todos";

    public async Task<List<TodoTask>> GetTasksAsync()
    {
        if (_appState.ActiveWorkspaceId == 0) return new List<TodoTask>();
        return await _http.GetFromJsonAsync<List<TodoTask>>(BaseUrl) ?? new List<TodoTask>();
    }

    public async Task AddTaskAsync(string title)
    {
        if (_appState.ActiveWorkspaceId == 0) return;
        var newTask = new { title = title };
        await _http.PostAsJsonAsync(BaseUrl, newTask);
    }

    public async Task DeleteTaskAsync(long id)
    {
        if (_appState.ActiveWorkspaceId == 0) return;
        await _http.DeleteAsync($"{BaseUrl}/{id}");
    }

    public async Task ToggleTaskAsync(long id)
    {
        if (_appState.ActiveWorkspaceId == 0) return;
        await _http.PatchAsync($"{BaseUrl}/{id}", null);
    }

    public async Task UpdateTaskTitleAsync(long id, string newTitle)
    {
        if (_appState.ActiveWorkspaceId == 0) return;
        var update = new { title = newTitle };
        await _http.PutAsJsonAsync($"{BaseUrl}/{id}", update);
    }
}