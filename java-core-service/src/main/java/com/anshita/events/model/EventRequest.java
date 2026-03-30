package com.anshita.events.model;

import jakarta.validation.constraints.NotBlank;
import jakarta.validation.constraints.NotNull;

import java.time.Instant;
import java.util.Map;
import java.util.UUID;

public record EventRequest(
        String eventId,
        @NotBlank String eventType,
        @NotBlank String partitionKey,
        @NotNull Map<String, Object> payload,
        int retryCount,
        Instant createdAt
) {
    public EventRequest normalize() {
        return new EventRequest(
                eventId == null || eventId.isBlank() ? UUID.randomUUID().toString() : eventId,
                eventType,
                partitionKey,
                payload,
                retryCount,
                createdAt == null ? Instant.now() : createdAt
        );
    }
}
