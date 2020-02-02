#version 410

uniform int renderType;
uniform int gameState;
uniform sampler2D diffuse;
uniform vec3 envLight;
uniform vec4 debugColor;

in vec2 fragTexCoord;
in vec3 fragNormal;

out vec4 fragColor;

void renderBackgroundSolid() {
  vec4 diffuseColor = texture2D(diffuse, fragTexCoord.st);
  fragColor = vec4(diffuseColor.rgb, 1.0);
  // Must override for all render functions
  gl_FragDepth = gl_FragCoord.z;
}

void renderBackgroundTransparent() {
  vec4 diffuseColor = texture2D(diffuse, fragTexCoord.st);
  fragColor = vec4(diffuseColor.rgb, diffuseColor.a);
  if (fragColor.a == 0) {
    // Hide transparent pixels
    gl_FragDepth = 1;
  } else {
    gl_FragDepth = gl_FragCoord.z;
  }
}

void renderCameraMask() {
  vec4 diffuseColor = texture2D(diffuse, fragTexCoord.st);
  fragColor = diffuseColor;

  if (fragColor.a == 0) {
    // Hide transparent pixels
    gl_FragDepth = 1;
  } else {
    gl_FragDepth = gl_FragCoord.z;
  }
}

void renderEntity() {
  vec4 diffuseColor = texture2D(diffuse, fragTexCoord.st);
  vec3 lightColor = envLight;
  fragColor = vec4(vec3(diffuseColor) * lightColor, 1.0);
  gl_FragDepth = gl_FragCoord.z;
}

void renderSprite() {
  vec4 diffuseColor = texture2D(diffuse, fragTexCoord.st);
  fragColor = diffuseColor;

  if (fragColor.a == 0) {
    // Hide transparent pixels
    gl_FragDepth = 1;
  } else {
    gl_FragDepth = gl_FragCoord.z;
  }
}

void renderItem() {
  vec4 diffuseColor = texture2D(diffuse, fragTexCoord.st);
  fragColor = vec4(diffuseColor.rgb, 1.0);
  gl_FragDepth = gl_FragCoord.z;
}

void renderDebug() {
  fragColor = debugColor;
  gl_FragDepth = gl_FragCoord.z;
}

void renderMainGame() {
  switch (renderType) {
    case -1:
      renderDebug();
      break;
    case 1:
      renderBackgroundSolid();
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
      renderBackgroundSolid();
      break;
    case 2:
      renderBackgroundTransparent();
      break;
  }
}
