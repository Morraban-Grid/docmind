# Use an official Python runtime as a parent image
FROM python:3.12.3

# Set the working directory in the container
WORKDIR /app

# Copy dependencies file
COPY requirements.txt .

# Install dependencies
RUN pip install --no-cache-dir -r requirements.txt

# Copy the rest of the application code
COPY . .

# Expose the port the app runs on
EXPOSE 8000

# Command to run the application
CMD ["python", "main.py"]