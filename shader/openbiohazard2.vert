#version 410

layout (location = 0) in vec3 position;
layout (location = 1) in vec2 vertTexCoord;
layout (location = 2) in vec3 vertNormal;

uniform int renderType;
uniform int gameState;
uniform mat4 projection;
uniform mat4 view;
uniform mat4 model;
// animation
uniform mat4 boneOffset;

out vec2 fragTexCoord;
out vec3 fragNormal;

void renderBackground2D() {
  gl_Position = vec4(position, 1.0);
  fragTexCoord = vertTexCoord;
}

void renderEntity() {
  vec4 modelPos = model * boneOffset * vec4(position, 1.0);
  gl_Position = projection * view * modelPos;
  fragTexCoord = vertTexCoord;

  fragNormal = mat3(transpose(inverse(model))) * vertNormal;
}

void renderSprite() {
  gl_Position = projection * view * vec4(position, 1.0);
  fragTexCoord = vertTexCoord;
}

void renderDebug() {
  gl_Position = projection * view * vec4(position, 1.0);
  fragTexCoord = vertTexCoord;
}

void renderItem() {
  gl_Position = projection * view * model * vec4(position, 1.0);
  fragTexCoord = vertTexCoord;
}

void renderMainGame() {
  switch (renderType) {
    case -1:
      renderDebug();
      break;
    case 1:
      renderBackground2D();
      break;
    case 2:
      renderBackground2D();
      break;
    case 3:
      renderEntity();
      break;
    case 4:
      renderSprite();
      break;
    case 5:
      renderItem();
      break;
  }
}

void main() {
  switch (gameState) {
    case 0:
      renderMainGame();
      break;
    case 1:
    case 2:
      renderBackground2D();
      break;
  }
}
