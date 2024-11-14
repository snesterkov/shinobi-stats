#!/usr/bin/python3

import subprocess
import json
import pymysql
import sys


connection = pymysql.connect(host='localhost', user='user', password='readpassword', database='ccio', cursorclass=pymysql.cursors.DictCursor)

monitors = {}

with connection:

    with connection.cursor() as cursor:
#        sql = "SELECT `mid`,`ke`,`name`,`protocol`,`host`,`path`,`port` from `Monitors`"
        sql = "SELECT `mid`,`ke`,`name`, json_extract(details, '$.auto_host') AS auto_host from `Monitors`"
        cursor.execute(sql)
        results = cursor.fetchall()
        for result in results:
#               monitors[result["mid"]]={"name": result["name"], "stream": f'{result["protocol"]}://{result["host"]}:{result["port"]}{result["path"]}'}
                stream = result["auto_host"].replace("\"", "")
                monitors[result["mid"]]={"name": result["name"], "stream": stream}
# Use public RTSP Stream for testing
probe_full = {'monitors': []}
for monitor in monitors:
        in_stream =  monitors[monitor]['stream']

        probe_command = ['/usr/local/bin/ffprobe',
                 '-v', 'panic',
#                 '-rtsp_transport', 'tcp',  # Force TCP (for testing)]
#                 '-select_streams', 'v:0',  # Select only video stream 0.
                 '-show_entries', 'stream', # Select only width and height entries
                 '-print_format', 'json', # Get output in JSON format
                 '-timeout', '3000000',
                 in_stream]

# Read video width, height using FFprobe:
        p0 = subprocess.Popen(probe_command, stdout=subprocess.PIPE)
        probe_str = p0.communicate()[0] # Reading content of p0.stdout (output of FFprobe) as string
        probe_json = probe_str.decode('utf8').replace("'", '"')
        p0.wait()
        probe_dict = json.loads(probe_str) # Convert string from JSON format to dictonary.

        if p0.returncode:
                status = 0
        else:
                status = 1
        probe_dict['name'] = monitors[monitor]['name']
        probe_dict['stream'] = monitors[monitor]['stream']
        probe_dict['status'] = status
        probe_dict['mid'] = monitor
        probe_full['monitors'].append(probe_dict)
probe_json = json.dumps(probe_full, indent=4, ensure_ascii=False)

f = open('/home/www/export_data.json', 'w')
f.write(probe_json)
f.close()