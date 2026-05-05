using System.Net.Http.Json;
using System.Threading.Tasks;
using System.Collections.Generic;
using note_nexus_client.Models;

namespace note_nexus_client.Services;

public class NoteService
{
    private readonly HttpClient _http;
    private readonly AppState _appState;

    public NoteService(HttpClient http, AppState appState)
    {
        _http = http;
        _appState = appState;
    }

    private string BaseUrl => $"v1/workspaces/{_appState.ActiveWorkspaceId}/notes";

    public async Task<List<Note>> GetNotesAsync()
    {
        if (_appState.ActiveWorkspaceId == 0) return new List<Note>();
        return await _http.GetFromJsonAsync<List<Note>>(BaseUrl) ?? new List<Note>();
    }

    public async Task CreateNoteAsync(string content)
    {
        if (_appState.ActiveWorkspaceId == 0) return;
        var newNote = new { content = content };
        await _http.PostAsJsonAsync(BaseUrl, newNote);
    }

    public async Task UpdateNoteAsync(long id, string content)
    {
        if (_appState.ActiveWorkspaceId == 0) return;
        var update = new { content = content };
        await _http.PutAsJsonAsync($"{BaseUrl}/{id}", update);
    }

    public async Task DeleteNoteAsync(long id)
    {
        if (_appState.ActiveWorkspaceId == 0) return;
        await _http.DeleteAsync($"{BaseUrl}/{id}");
    }
}