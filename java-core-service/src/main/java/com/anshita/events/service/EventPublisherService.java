package com.anshita.events.service;

import com.anshita.events.model.EventRequest;
import io.micrometer.core.instrument.Counter;
import io.micrometer.core.instrument.MeterRegistry;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.kafka.core.KafkaTemplate;
import org.springframework.stereotype.Service;

import java.util.List;
import java.util.concurrent.CompletableFuture;

@Service
public class EventPublisherService {

    private final KafkaTemplate<String, EventRequest> kafkaTemplate;
    private final String mainTopic;
    private final Counter publishedCounter;

    public EventPublisherService(
            KafkaTemplate<String, EventRequest> kafkaTemplate,
            @Value("${app.topics.main}") String mainTopic,
            MeterRegistry meterRegistry
    ) {
        this.kafkaTemplate = kafkaTemplate;
        this.mainTopic = mainTopic;
        this.publishedCounter = meterRegistry.counter("events_published_total");
    }

    public CompletableFuture<Void> publish(EventRequest request) {
        EventRequest normalized = request.normalize();
        return kafkaTemplate.send(mainTopic, normalized.partitionKey(), normalized)
                .thenAccept(result -> publishedCounter.increment());
    }

    public void publishBatch(List<EventRequest> requests) {
        requests.forEach(this::publish);
    }
}
