#version 330 core
out vec4 FragColor;

in vec2 TexCoord;

uniform sampler2D aTexture;
void main()
{
	FragColor = texture(aTexture, TexCoord);
}
