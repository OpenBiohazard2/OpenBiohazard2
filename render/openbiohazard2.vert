#version 410

layout (location = 0) in vec3 position;
layout (location = 1) in vec2 vertTexCoord;
layout (location = 2) in vec3 vertNormal;

uniform int renderType;
uniform mat4 projection;
uniform mat4 view;
uniform mat4 model;
// animation
uniform mat4 boneOffset;

out vec2 fragTexCoord;
out vec3 fragNormal;

void renderBackground() {
  gl_Position = vec4(position, 1.0);
  fragTexCoord = vertTexCoord;
}

void renderCameraMask() {
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

void main() {
	switch (renderType) {
    case 1:
      renderBackground();
      break;
    case 2:
      renderCameraMask();
      break;
    case 3:
      renderEntity();
      break;
    case 4:
      renderSprite();
      break;
    case 5:
      renderDebug();
      break;
  }
}
