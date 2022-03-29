FROM python:3.9
WORKDIR /src
COPY . .
RUN pip install -r requirements.txt

EXPOSE 5001
ENTRYPOINT ["python", "main.py"]