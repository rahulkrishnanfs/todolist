#! /bin/bash
docker run -p 8080:8080 -v "$(pwd)/secrets:/etc/todolist/secrets:ro" rahulkrishnanfs/todolist:v1