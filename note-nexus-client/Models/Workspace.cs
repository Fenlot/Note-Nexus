using System;

namespace note_nexus_client.Models
{
    public class Workspace
    {
        public long Id { get; set; }
        public string Name { get; set; } = string.Empty;
        public long OwnerId { get; set; }
        public string SubscriptionTier { get; set; } = "free";
        public DateTime CreatedAt { get; set; }
    }
}
