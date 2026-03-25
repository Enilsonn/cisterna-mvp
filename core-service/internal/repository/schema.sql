CREATE EXTENSION IF NOT EXISTS postgis;

CREATE TABLE IF NOT EXISTS truck_coordinates (
    id SERIAL,
    truck_id VARCHAR(50) NOT NULL,
    location GEOMETRY(Point, 4326) NOT NULL,
    recorded_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id, recorded_at) 
) PARTITION BY RANGE (recorded_at);

DO $$
DECLARE
    data_atual DATE := '2026-01-01';
    data_limite DATE := '2030-12-01';
    nome_tabela TEXT;
BEGIN
    WHILE data_atual <= data_limite LOOP
        -- truck_coordinates_2026_03 (exemplo)
        nome_tabela := 'truck_coordinates_' || to_char(data_atual, 'YYYY_MM');
        
        -- o postgres soma 1 mês automaticamente (ele lida com a virada de ano sozinho)
        EXECUTE format('
            CREATE TABLE IF NOT EXISTS %I PARTITION OF truck_coordinates
            FOR VALUES FROM (%L) TO (%L);',
            nome_tabela,
            data_atual,
            data_atual + INTERVAL '1 month'
        );
        
        data_atual := data_atual + INTERVAL '1 month';
    END LOOP;
END $$;

CREATE INDEX IF NOT EXISTS idx_truck_id ON truck_coordinates(truck_id);
CREATE INDEX IF NOT EXISTS idx_recorded_at ON truck_coordinates(recorded_at);
CREATE INDEX IF NOT EXISTS idx_location_gist ON truck_coordinates USING GIST(location);