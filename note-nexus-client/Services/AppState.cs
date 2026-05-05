using System;

namespace note_nexus_client.Services
{
    public class AppState
    {
        public long ActiveWorkspaceId { get; private set; }
        
        public event Action? OnChange;

        public void SetActiveWorkspace(long workspaceId)
        {
            if (ActiveWorkspaceId != workspaceId)
            {
                ActiveWorkspaceId = workspaceId;
                NotifyStateChanged();
            }
        }

        private void NotifyStateChanged() => OnChange?.Invoke();
    }
}
