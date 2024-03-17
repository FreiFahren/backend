# flake8: noqa

import json
import fuzzywuzzy.process

with open('data/lines.json') as f:
    LinesData = json.load(f)

with open('data/stations.json') as f:
    StationsData = json.load(f)

length = 0
for lines in LinesData:
    for _ in LinesData[lines]:
        length += 1

count = 0

choices = []

for id in StationsData:
    choices.append(StationsData[id]['name'])

for lines in LinesData:
    for station in range(0, len(LinesData[lines])):
        station_name = LinesData[lines][station]
        closest_match = fuzzywuzzy.process.extractOne(station_name, choices)
    
        if closest_match:
            count += 1

            for id in StationsData:
                if StationsData[id]['name'] == closest_match[0]:
                    closest_station_id = id

            
            LinesData[lines][station] = closest_station_id
            print(f"{count}/{length} - {station_name} -> {closest_station_id} ({closest_match[1]}%)")

# Dump the LinesData into lineslist.json
with open('data/lineslist.json', 'w') as f:
    json.dump(LinesData, f, ensure_ascii=False, indent=4)
