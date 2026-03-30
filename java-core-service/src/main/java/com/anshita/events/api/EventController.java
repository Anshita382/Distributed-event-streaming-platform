package com.anshita.events.api;

import com.anshita.events.model.EventRequest;
import com.anshita.events.service.EventPublisherService;
import jakarta.validation.Valid;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.bind.annotation.RestController;

import java.time.Instant;
import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.UUID;

@RestController
@RequestMapping("/api/events")
public class EventController {

    private final EventPublisherService publisherService;

    public EventController(EventPublisherService publisherService) {
        this.publisherService = publisherService;
    }

    @PostMapping
    public ResponseEntity<Map<String, Object>> publish(@Valid @RequestBody EventRequest request) {
        publisherService.publish(request);
        return ResponseEntity.accepted().body(Map.of(
                "status", "queued",
                "eventId", request.normalize().eventId()
        ));
    }

    @PostMapping("/load-test")
    public ResponseEntity<Map<String, Object>> loadTest(@RequestParam(defaultValue = "1000") int count) {
        List<EventRequest> events = new ArrayList<>(count);
        for (int i = 0; i < count; i++) {
            Map<String, Object> payload = new HashMap<>();
            payload.put("sequence", i);
            payload.put("source", "load-test-api");
            payload.put("amount", i % 250);
            payload.put("createdAt", Instant.now().toString());
            payload.put("simulateFailure", i % 200 == 0);

            events.add(new EventRequest(
                    UUID.randomUUID().toString(),
                    i % 2 == 0 ? "order.created" : "payment.captured",
                    "customer-" + (i % 500),
                    payload,
                    0,
                    Instant.now()
            ));
        }
        publisherService.publishBatch(events);
        return ResponseEntity.accepted().body(Map.of(
                "status", "queued",
                "count", count
        ));
    }

    @GetMapping("/health")
    public ResponseEntity<Map<String, String>> health() {
        return ResponseEntity.ok(Map.of("status", "up"));
    }
}
