using System.Net.Http;
using System.Net.Http.Json;
using System.Threading.Tasks;
using Microsoft.AspNetCore.Components.Authorization;
using note_nexus_client.Models;

namespace note_nexus_client.Services
{
    public class AuthService
    {
        private readonly HttpClient _http;
        private readonly AuthenticationStateProvider _authStateProvider;
        private readonly LocalStorageService _localStorage;
        private const string TokenKey = "authToken";

        public AuthService(HttpClient http, AuthenticationStateProvider authStateProvider, LocalStorageService localStorage)
        {
            _http = http;
            _authStateProvider = authStateProvider;
            _localStorage = localStorage;
        }

        public async Task<bool> LoginAsync(LoginRequest request)
        {
            var response = await _http.PostAsJsonAsync("v1/login", request);

            if (!response.IsSuccessStatusCode)
            {
                return false;
            }

            var authResponse = await response.Content.ReadFromJsonAsync<AuthResponse>();
            if (authResponse != null && !string.IsNullOrWhiteSpace(authResponse.Token))
            {
                await _localStorage.SetItemAsync(TokenKey, authResponse.Token);
                ((CustomAuthStateProvider)_authStateProvider).MarkUserAsAuthenticated(authResponse.Token);
                return true;
            }

            return false;
        }

        public async Task<bool> SignupAsync(SignupRequest request)
        {
            var response = await _http.PostAsJsonAsync("v1/signup", request);

            if (!response.IsSuccessStatusCode)
            {
                return false;
            }

            var authResponse = await response.Content.ReadFromJsonAsync<AuthResponse>();
            if (authResponse != null && !string.IsNullOrWhiteSpace(authResponse.Token))
            {
                await _localStorage.SetItemAsync(TokenKey, authResponse.Token);
                ((CustomAuthStateProvider)_authStateProvider).MarkUserAsAuthenticated(authResponse.Token);
                return true;
            }

            return false;
        }

        public async Task LogoutAsync()
        {
            await _localStorage.RemoveItemAsync(TokenKey);
            ((CustomAuthStateProvider)_authStateProvider).MarkUserAsLoggedOut();
        }
    }
}
