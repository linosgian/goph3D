#version 330 core
layout(location = 0) in vec2 position;
layout(location = 1) in vec3 aColor;
layout(location = 2) in vec2 aTexCoord;

out vec2 TexCoord;
out vec3 ourColor;

uniform mat4 model;
uniform mat4 projection;
uniform mat4 view;

void main()
{
	gl_Position = projection * view * model * vec4(position, 0.0, 1.0);
	TexCoord = aTexCoord;
	ourColor = aColor;
}
