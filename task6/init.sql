CREATE DATABASE telemetry;
CREATE DATABASE smarthome;
CREATE DATABASE devices;

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

-- Connect to the database
\c smarthome;

-- Create the sensors table
CREATE TABLE IF NOT EXISTS sensors (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    type VARCHAR(50) NOT NULL,
    location VARCHAR(100) NOT NULL,
    value FLOAT DEFAULT 0,
    unit VARCHAR(20),
    status VARCHAR(20) NOT NULL DEFAULT 'inactive',
    last_updated TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Create indexes for common queries
CREATE INDEX IF NOT EXISTS idx_sensors_type ON sensors(type);
CREATE INDEX IF NOT EXISTS idx_sensors_location ON sensors(location);
CREATE INDEX IF NOT EXISTS idx_sensors_status ON sensors(status);

\c devices;

-- Создание таблицы Sensors
CREATE TABLE IF NOT EXISTS Sensors (
    Id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    Name VARCHAR(100) NOT NULL,
    Type VARCHAR(50) NOT NULL,
    Location VARCHAR(200) NULL,
    IsActive BOOLEAN NOT NULL DEFAULT true,
    CreatedAt TIMESTAMP NOT NULL DEFAULT NOW(),
    UpdatedAt TIMESTAMP NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS IX_Sensors_Name ON Sensors (Name);

CREATE INDEX IF NOT EXISTS IX_Sensors_IsActive ON Sensors (IsActive);

CREATE INDEX IF NOT EXISTS IX_Sensors_Type ON Sensors (Type);

INSERT INTO Sensors (Id, Name, Type, Location, IsActive) VALUES
    ('550e8400-e29b-41d4-a716-446655440000', 'Temperature Living Room', 'Thermometer', 'Living Room', true),
    ('550e8400-e29b-41d4-a716-446655440001', 'Motion Entrance', 'Motion Sensor', 'Entrance', true),
    ('550e8400-e29b-41d4-a716-446655440002', 'Humidity Kitchen', 'Hygrometer', 'Kitchen', true),
    ('550e8400-e29b-41d4-a716-446655440003', 'Light Bedroom', 'Light Sensor', 'Bedroom', false)
ON CONFLICT (Name) DO NOTHING;