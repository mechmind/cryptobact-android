#include <GLES2/gl2.h>

void set_ortho_proj(GLfloat *matrix, GLfloat left, GLfloat right,
        GLfloat bottom, GLfloat top, GLfloat near, GLfloat far)
{
  matrix[0] = 2.0f / (right - left);
  matrix[5] = 2.0f / (top - bottom);
  matrix[10] = -2.0 / (far - near);
  matrix[1] = matrix[2] = matrix[3] = 0.0f;
  matrix[4] = matrix[6] = matrix[7] = 0.0f;
  matrix[8] = matrix[9] = matrix[11] = 0.0f;
  matrix[12] = -(right + left) / (right - left);
  matrix[13] = -(top + bottom) / (top - bottom);
  matrix[14] = -(far + near) / (far - near);
  matrix[15] = 1.0f;
};
