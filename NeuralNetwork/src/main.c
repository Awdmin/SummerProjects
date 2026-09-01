#include <stdio.h>
#include <stdlib.h>

#include "nn.h"


int main() {

    Layer* l1 = init_layer(3, 3);
    Layer* l2 = init_layer(3, 4);
    Layer* l3 = init_layer(4, 2);

    Matrix* v = (Matrix*)malloc(sizeof(Matrix));
    v->cols = 1;
    v->rows = 3;
    v->v = (float**)malloc(sizeof(float*)*v->cols);
    v->v[0] = (float*)malloc(sizeof(float));
    v->v[0][0] = 1;
    v->v[1] = (float*)malloc(sizeof(float));
    v->v[1][0] = 2;
    v->v[2] = (float*)malloc(sizeof(float));
    v->v[2][0] = 3;


    Network* nn = (Network*)malloc(sizeof(Network));
    nn->n_layers = 3;
    Layer* layers = (Layer*)malloc(sizeof(Layer)*nn->n_layers);
    layers[0].biases = l1->biases;
    layers[0].weights = l1->weights;
    layers[1].biases = l2->biases;
    layers[1].weights = l2->weights;
    layers[2].biases = l3->biases;
    layers[2].weights = l3->weights;
    nn->layers = layers;


    Matrix* out = forward_pass(nn, v);

    print_matrix(out);

    return 0;
}
