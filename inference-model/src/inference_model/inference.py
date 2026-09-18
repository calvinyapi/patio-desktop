import cv2 as cv


cap = cv.VideoCapture(0)

def inference():
    ret, frame = cap.read()
    if not ret:
        print("Failed to capture frame from camera.")
    