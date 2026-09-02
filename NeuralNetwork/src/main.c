#include <stdio.h>
#include <stdlib.h>

#include "nn.h"


int main() {

    Matrix* v = (Matrix*)malloc(sizeof(Matrix));
    v->cols = 1;
    v->rows = 3;
    v->v = (double**)malloc(sizeof(double*)*v->cols);
    v->v[0] = (double*)malloc(sizeof(double));
    v->v[0][0] = 1;
    v->v[1] = (double*)malloc(sizeof(double));
    v->v[1][0] = 2;
    v->v[2] = (double*)malloc(sizeof(double));
    v->v[2][0] = 3;


    Network* nn = (Network*)malloc(sizeof(Network));
    nn->n_layers = 3;
    Layer* layers = (Layer*)malloc(sizeof(Layer)*nn->n_layers);
    layers[0] = *init_layer(3, 32, 0);
    layers[1] = *init_layer(32, 32, 0);
    layers[2] = *init_layer(32, 4, 3);
    nn->layers = layers;

    print_network(nn);


    Matrix** r = (Matrix**)malloc(sizeof(Matrix*)*nn->n_layers);

    Matrix* out = forward_pass(nn, v, r);

    print_matrix(out);

    return 0;
}
