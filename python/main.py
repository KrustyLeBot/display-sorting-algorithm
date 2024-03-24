import sys
import pygame as pg
import numpy as np
import copy
from threading import Thread

class Display:
    def __init__(self, sorter):
        pg.init()
        self.screen = pg.display.set_mode((1024,1024))
        self.clock = pg.time.Clock()
        self.sorter = sorter

    def run(self):
        while True:
            self.check_events()
            self.update()
            self.draw()

    def check_events(self):
        events = pg.event.get()
        for event in events:
            if event.type == pg.QUIT or (event.type == pg.KEYDOWN and event.key == pg.K_ESCAPE):
                sys.exit()

    def update(self):
        self.clock.tick(30)
        pg.display.set_caption(f"fps: {self.clock.get_fps() :.1f}")

    def draw(self):
        # reshape into 2d array and remove weight to render it
        tmp = self.sorter.cells.reshape((-1, 1024, 4))[:, :, :3]
        surf = pg.Surface((tmp.shape[0], tmp.shape[1]))
        pg.surfarray.blit_array(surf, tmp)
        surf = pg.transform.scale(surf, (1024, 1024))

        self.screen.blit(surf, (0, 0))
        pg.display.update()


class Sorter:
    def __init__(self):

        image = pg.image.load("example.jpg")
        pixel = pg.PixelArray(image)

        self.cells = np.ndarray((1024, 1024, 4))

        weight = 0
        for i in range(self.cells.shape[0]):
            for j in range(self.cells.shape[1]):
                color = image.unmap_rgb(pixel[i, j])
                self.cells[i][j] = (color.r, color.g, color.b, weight)
                weight += 1

        # randomize pixels and keep 1d array of colors + weight
        self.cells = self.cells.reshape(-1, self.cells.shape[-1])
        np.random.shuffle(self.cells)
        np.sort(self.cells, axis=0)

        ######################################################################################################
        # replace sorting method to use here
        # start sort in another thread
        self.cells_progress = copy.deepcopy(self.cells)
        thread = Thread(target = self.quicksort, args = (self.cells_progress, 0, len(self.cells_progress)))
        thread.start()
        ######################################################################################################

    def quicksort(self, arr, start, stop):
        if len(arr) <= 1:
            return arr
            
        pivot = arr[len(arr) // 2]
            
        left = [x for x in arr if x[3] < pivot[3]]
        len_l = len(left)
            
        middle = [x for x in arr if x[3] == pivot[3]]
        start_m = start + len_l
        len_m = len(middle)
            
        right = [x for x in arr if x[3] > pivot[3]]
        start_r = start + len_l + len_m
        len_r = len(right)

        for i in range(0, len_l):
            self.cells[start + i] = left[i]
        for i in range(0, len_m):
            self.cells[start_m + i] = middle[i]
        for i in range(0, len_r):
            self.cells[start_r + i] = right[i]

        return self.quicksort(left, start, start + len_l) + middle + self.quicksort(right, start_r, stop)

if __name__ == "__main__":
    sorter = Sorter()
    display = Display(sorter)
    display.run()