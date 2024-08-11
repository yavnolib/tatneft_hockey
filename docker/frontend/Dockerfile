FROM python:3.10-slim

WORKDIR /app

# Dependencies
COPY tathockey/frontend/requirements.txt .
RUN python -m pip install -r requirements.txt

RUN mkdir -p /app/tathockey/backend/retinaface && \
    mkdir -p /app/tathockey/frontend

COPY tathockey/frontend /app/tathockey/frontend/

COPY tathockey/processor_utils.py /app/tathockey/

COPY setup.py .

RUN python -m pip install -e .

EXPOSE 8080

CMD ["python", "/app/tathockey/frontend/run.py"]
