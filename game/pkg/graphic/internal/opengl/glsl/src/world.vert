#version 330 core

layout (location = 0) in vec2 a_pos;

uniform mat4 u_model;
uniform mat4 u_view;

layout (std140) uniform ProjView {
    mat4 u_proj;
};

void main() {
    gl_Position = u_proj * u_view * u_model * vec4(a_pos, 0.0, 1.0);
}
