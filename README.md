
```
go clone git@github.com:mmirolim/maze.git
cd maze
go build

./maze
  -h int
    	height of maze (default 10)
  -w int
    	width of maze (default 10)
Gen Time 180.125µs
+---+---+---+---+---+---+---+---+---+---+
| S         |                           |
+---+---+   +---+---+   +---+---+---+   +
|       |               |   |       |   |
+   +---+---+---+---+---+   +   +   +   +
|               |               |   |   |
+   +---+---+   +   +---+---+---+   +   +
|   |       |   |   |           |   |   |
+   +   +   +   +   +   +---+   +   +   +
|   |   |   |       |       |       |   |
+   +   +   +---+---+---+   +---+---+   +
|       |   |   |       |           |   |
+---+---+   +   +   +   +---+---+   +   +
|       |   |       |       |   |       |
+   +   +   +   +---+   +   +   +---+---+
|   |       |   |       |       |       |
+   +---+---+---+   +---+---+   +---+   +
|       |           |       |       |   |
+---+   +---+   +---+   +---+---+   +   +
|               |                     E |
+---+---+---+---+---+---+---+---+---+---+

DrawTxt Time 380.5µs
Draw Maze Time 5.993221678s
Draw Maze Gen Time 6.001344901s
PathFound? with Recursion true cells visited 4669, path len 1831 Time 1.151957ms
Draw Maze Find Path Time 2.908757757s
ShortestPathFound? with BFS true cells visited 2796 path len 1831 Time 1.234331ms
Draw Maze Find Shortest Path Time 5.461315859s
```

![Maze](https://image.ibb.co/jR3NN9/maze.gif)
![Maze anim path with dfs recur](https://image.ibb.co/eMGL9p/maze_anim_path_with_dfs_recur.gif)
![Maze anim gen with dfs](https://image.ibb.co/m9b929/maze_anim_gen_with_dfs.gif)
![Maze anim path with bfs weight](https://image.ibb.co/cdXSpp/maze_anim_path_with_bfs_weight.gif)

