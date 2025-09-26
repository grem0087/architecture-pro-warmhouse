CREATE DATABASE telemetry;

\c telemetry;

CREATE TABLE IF NOT EXISTS sensors (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    location_id UUID NOT NULL,
    sensor_type VARCHAR(50) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS telemetry_data (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    sensor_id UUID NOT NULL REFERENCES sensors(id),
    value DECIMAL(10,4) NOT NULL,
    timestamp TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_telemetry_sensor_id ON telemetry_data(sensor_id);
CREATE INDEX IF NOT EXISTS idx_telemetry_timestamp ON telemetry_data(timestamp);
CREATE INDEX IF NOT EXISTS idx_telemetry_sensor_timestamp ON telemetry_data(sensor_id, timestamp);

INSERT INTO sensors (id, name, location_id, sensor_type) VALUES
    ('550e8400-e29b-41d4-a716-446655440000', 'Температура гостиная', '750e8400-e29b-41d4-a716-446655440001', 'temperature'),
    ('550e8400-e29b-41d4-a716-446655440002', 'Влажность кухня', '750e8400-e29b-41d4-a716-446655440001', 'humidity'),
    ('550e8400-e29b-41d4-a716-446655440003', 'Освещенность спальня', '750e8400-e29b-41d4-a716-446655440002', 'light')
ON CONFLICT (id) DO NOTHING;

INSERT INTO telemetry_data (sensor_id, value, timestamp) VALUES
    ('550e8400-e29b-41d4-a716-446655440000', 22.5, NOW() - INTERVAL '1 hour'),
    ('550e8400-e29b-41d4-a716-446655440000', 23.1, NOW() - INTERVAL '30 minutes'),
    ('550e8400-e29b-41d4-a716-446655440002', 45.0, NOW() - INTERVAL '2 hours'),
    ('550e8400-e29b-41d4-a716-446655440002', 47.2, NOW() - INTERVAL '1 hour'),
    ('550e8400-e29b-41d4-a716-446655440003', 350.0, NOW() - INTERVAL '3 hours');