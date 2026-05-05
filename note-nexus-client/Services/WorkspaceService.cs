using System.Collections.Generic;
using System.Net.Http;
using System.Net.Http.Json;
using System.Threading.Tasks;
using note_nexus_client.Models;

namespace note_nexus_client.Services
{
    public class WorkspaceService
    {
        private readonly HttpClient _http;

        public WorkspaceService(HttpClient http)
        {
            _http = http;
        }

        public async Task<List<Workspace>> GetWorkspacesAsync()
        {
            return await _http.GetFromJsonAsync<List<Workspace>>("v1/workspaces") ?? new List<Workspace>();
        }
    }
}
