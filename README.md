# OpenBiohazard2
Open source re-implementation of the original Resident Evil 2 engine written in Go and OpenGL. You must own a copy of the original game.

<div style="display:inline-block;">
<img src="https://github.com/OpenBiohazard2/OpenBiohazard2/raw/master/screenshots/beginning.png" alt="beginning" width="400" height="300" />
<img src="https://github.com/OpenBiohazard2/OpenBiohazard2/raw/master/screenshots/inventory.png" alt="inventory" width="400" height="300" />
</div>

### Installation

1. Clone this project.
2. Get the game data from your installed location. Copy all the files to the `data/` folder in this repository.
3. Run `go build`.

### Task list

- [ ] Audio
  - [ ] Background music
  - [ ] Core sound
  - (Audio file formats are parsed; playback is not implemented yet.)
- [ ] Game
  - [x] Collision detection
  - [x] Event triggers
  - [ ] Inventory system (UI/menu rendering implemented; item use and management pending)
  - [ ] Enemy AI (enemies spawn and render as debug placeholders; no AI or full models yet)
  - [ ] Puzzles
  - [ ] Door transitions (room switching works; door models and animations pending)
- [ ] Renderer
  - [x] Animation
  - [x] Pre-rendered background
  - [x] Depth testing
  - [ ] Sprites (assets are parsed; rendering is not enabled yet)
  - [ ] Shadows
- [ ] Scripting engine
  - [x] Control flow (if/else, loops, switch, goto/gosub)
  - [x] Camera changes
  - [x] AOT triggers, items, and doors
  - [ ] Full opcode coverage (44 of 103 opcodes implemented)

### Controls

- W/S to move forward/backward.
- A/D to rotate left/right.
- Tab to access inventory.
- Enter is action button.

### Other tools

- [Bio2 script viewer](https://github.com/OpenBiohazard2/Bio2ScriptIde)

### References

- https://github.com/pmandin
- https://github.com/yanmingsohu/Biohazard2
- https://github.com/MeganGrass/BioScript
- https://github.com/mortician
