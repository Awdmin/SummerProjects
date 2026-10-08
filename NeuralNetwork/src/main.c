#include <stdio.h>
#include <stdlib.h>

#include "nn.h"

const double LR = 1e-4;


int main() {

    Matrix* v = (Matrix*)malloc(sizeof(Matrix));
    v->cols = 1;
    v->rows = 3;
    v->v = (double**)malloc(sizeof(double*)*v->cols);
    v->v[0] = (double*)malloc(sizeof(double));
    v->v[0][0] = _rand();
    v->v[1] = (double*)malloc(sizeof(double));
    v->v[1][0] = _rand();
    v->v[2] = (double*)malloc(sizeof(double));
    v->v[2][0] = _rand();

    Matrix* t = (Matrix*)malloc(sizeof(Matrix));
    t->cols = 1;
    t->rows = 4;
    t->v = (double**)malloc(sizeof(double*)*v->cols);
    t->v[0] = (double*)malloc(sizeof(double));
    t->v[0][0] = _rand();
    t->v[1] = (double*)malloc(sizeof(double));
    t->v[1][0] = _rand();
    t->v[2] = (double*)malloc(sizeof(double));
    t->v[2][0] = _rand();
    t->v[3] = (double*)malloc(sizeof(double));
    t->v[3][0] = _rand();

    Network* nn = (Network*)malloc(sizeof(Network));
    nn->n_layers = 3;
    Layer* layers = (Layer*)malloc(sizeof(Layer)*nn->n_layers);
    layers[0] = *init_layer(3, 32, 0);
    layers[1] = *init_layer(32, 32, 0);
    layers[2] = *init_layer(32, 4, 3);
    nn->layers = layers;

    print_network(nn);


    Matrix** z = (Matrix**)malloc(sizeof(Matrix*)*nn->n_layers);
    Matrix** a = (Matrix**)malloc(sizeof(Matrix*)*nn->n_layers);

    print_matrix(v);

    Matrix* out = forward_pass(nn, v, z, a);

    print_matrix(out);

    print_layer(&nn->layers[2]);
    backward_pass(nn, a, z, t, LR);

    return 0;
}
