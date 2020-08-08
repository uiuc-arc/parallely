import numpy as np
import os
import keras
from keras.models import Sequential, load_model
from keras.layers import Dense, Dropout, Flatten, Activation
from keras.layers import Conv2D, MaxPooling2D, ZeroPadding2D
from keras.callbacks import ModelCheckpoint
from keras.optimizers import SGD
import cv2
import random

import matplotlib.pyplot as plt

from keras import backend as K
#from frontend.approxhpvm_translator import translate_to_approxhpvm
import pdb

batch_size = 32
num_classes = 5
img_size = 32

#run_dir = './cv/CNN_MIO_KERAS/'
run_dir = './'


def mio_model():
    model = Sequential()

    model.add(Conv2D(32, (3, 3), activation='relu', input_shape=(3, img_size, img_size)))
    model.add(Conv2D(32, (3, 3), activation='relu'))
    model.add(MaxPooling2D(pool_size=(2, 2)))
    model.add(Dropout(0.25))

    model.add(Conv2D(64, (3, 3), activation='relu'))
    model.add(Conv2D(64, (3, 3), activation='relu'))
    model.add(MaxPooling2D(pool_size=(2, 2)))
    model.add(Dropout(0.25))

    model.add(Flatten())
    model.add(Dense(256, activation='relu'))
    model.add(Dropout(0.5))
    model.add(Dense(num_classes, activation='softmax'))
    return model

def _get_images_labels(dsplit):
  if dsplit == 'train':  
    features = np.load(run_dir + 'training_images.npy')
    labels_t = np.load(run_dir + 'training_labels.npy')
    shuffle_buffer_size = 6000
  else:  
    features = np.load(run_dir + 'test_images.npy')
    labels_t = np.load(run_dir + 'test_labels.npy')
    shuffle_buffer_size = 1000
  assert features.shape[0] == labels_t.shape[0]
  labels = labels_t 
  return features,labels


def calculate_accuracy(conf, true_result, prediction, min_conf, max_conf):
  for i, line in enumerate(conf):
      max_weight.append(max(line))
      if max(line) == 1:
          max_weight_indices.append(i)    
  return 0


if __name__ == "__main__":   

  # Read the training data set
  K.set_image_data_format('channels_first')
  features,labels = _get_images_labels('train')
  test_features,test_labels = _get_images_labels('test')
  model = mio_model()
  model.load_weights(run_dir + 'model.h5')
  
  # Pick a random image from the test set and get its output
  random_image_number = random.randrange(len(test_labels))

  print(test_features.shape, test_labels.shape, test_features[0].shape, test_features[0,:].shape)
  

  outputs = model.predict(np.array([test_features[random_image_number],]), batch_size=1)

  # print(len(test_labels))
  # print(outputs)
  print("Prediction: ", outputs.argmax(1)[0], max(outputs[0]))
