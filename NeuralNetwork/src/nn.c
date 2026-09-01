#include <stdlib.h>
#include <stdio.h>

#include "nn.h"

/*
################### What is done so far ###################
- vector operations
- basic structs for neural network
- forward pass
 
################### TO DO ####################
- activation functions
- backwards pass (learning)
- optimizer

*/

float _rand() {
    return rand() / (float)(RAND_MAX/2) - 1;
}

Vector* sum(Vector* a, Vector* b) {

    if(a->dim != b->dim) {
        return NULL;
    }

    Vector* r = (Vector*)malloc(sizeof(Vector));
    r->dim = a->dim;
    float* v = (float*)malloc(sizeof(int)*r->dim);

    for(int i = 0; i < r->dim; i++) {
        v[i] = a->v[i] + b->v[i];
    }
    r->v = v;
    return r;
}


float dot(Vector* a, Vector* b) {

    if(a->dim != b->dim) {
        return 0;
    }

    float s = 0;
    for(int i = 0; i < a->dim; i++) {
        s += a->v[i]*b->v[i];
    }
    return s;
}

Vector* layer_out(Layer* l, Vector* v_in) {
    if(l->nodes[0].weights->dim != v_in->dim) {
        return NULL;
    }
    Vector* r = (Vector*)malloc(sizeof(Vector));
    r->dim = l->n_nodes;
    r->v = (float*)malloc(sizeof(float)*r->dim);
    for(int i = 0; i < l->n_nodes; i++) {
        r->v[i] = dot(l->nodes[i].weights, v_in) + l->nodes[i].bias;
    }

    return r;
}

Vector* forward_pass(Network* nn, Vector* input) {
    if(nn->n_layers == 0) {
        return NULL;
    }

    Vector* v = layer_out(&nn->layers[0], input);
    for(int i = 1; i < nn->n_layers; i++) {
        v = layer_out(&nn->layers[i], v);
    }
    return v;
}



// ----- Init functions -----

Layer* init_layer(int n, int dim) { //dim needs to be the same as the n_nodes on the previous layer

    Layer* l = (Layer*)malloc(sizeof(Layer));
    l->nodes = (Node*)malloc(sizeof(Node)*n);
    l->n_nodes = n;

    for(int i = 0; i < n; i++) {
        l->nodes[i].bias = _rand();
        l->nodes[i].weights = (Vector*)malloc(sizeof(Vector));
        l->nodes[i].weights->dim = dim;
        l->nodes[i].weights->v = (float*)malloc(sizeof(float)*dim);
        for(int k = 0; k < dim; k++) {
            l->nodes[i].weights->v[k] = _rand();
        }
    }

    return l;
}



// ----- Print functions -----

void print_vector(Vector* a) {
    printf("[");
    for(int i = 0; i < a->dim; i++) {
        printf(" %f", a->v[i]);
    }
    printf(" ]\n");
}

void print_layer(Layer* l) {
    for(int i = 0; i < l->n_nodes; i++) {
        printf("n%d_ bias: %f, weights: ", i, l->nodes[i].bias);
        print_vector(l->nodes[i].weights);
    }
}
