package com.anshita.events.config;

import org.apache.kafka.clients.admin.NewTopic;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.kafka.config.TopicBuilder;

@Configuration
public class KafkaTopicConfig {

    @Bean
    public NewTopic mainTopic(@Value("${app.topics.main}") String topic) {
        return TopicBuilder.name(topic).partitions(6).replicas(1).build();
    }

    @Bean
    public NewTopic retryTopic(@Value("${app.topics.retry}") String topic) {
        return TopicBuilder.name(topic).partitions(6).replicas(1).build();
    }

    @Bean
    public NewTopic dlqTopic(@Value("${app.topics.dlq}") String topic) {
        return TopicBuilder.name(topic).partitions(6).replicas(1).build();
    }
}
