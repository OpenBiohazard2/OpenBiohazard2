# OpenBiohazard2
Open source re-implementation of the original Resident Evil 2 engine written in Go and OpenGL. You must own a copy of the original game.

<div style="display:inline-block;">
<img src="https://github.com/samuelyuan/OpenBiohazard2/raw/master/screenshots/beginning.png" alt="beginning" width="400" height="300" />
<img src="https://github.com/samuelyuan/OpenBiohazard2/raw/master/screenshots/inventory.png" alt="inventory" width="400" height="300" />
</div>

### Installation

1. Clone this project.
2. Download the dependencies

```
go get github.com/go-gl/gl/v4.1-core/gl
go get github.com/go-gl/glfw/v3.2/glfw
go get github.com/go-gl/mathgl/mgl32
```

3. Get the game data and copy all the files to the `data/` folder in this repository.
4. Run `go build`.
