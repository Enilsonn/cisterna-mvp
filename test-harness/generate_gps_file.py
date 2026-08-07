import json
from pathlib import Path

out = Path("test-harness/gps_points.json")
points = []
base_lat = -7.12
base_lon = -34.88
for i in range(1000):
    points.append({
        "lat": round(base_lat + i * 0.0001, 6),
        "lon": round(base_lon + i * 0.0001, 6),
        "timestamp": f"2026-08-07T00:{i//60:02d}:{i%60:02d}Z"
    })

out.write_text(json.dumps(points, indent=2), encoding="utf-8")
print(f"generated {len(points)} points at {out}")
