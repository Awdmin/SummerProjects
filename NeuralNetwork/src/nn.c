#include <iso646.h>
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

Matrix* sum(Matrix* a, Matrix* b) {

    if(a->cols != b->cols || a->rows != b->rows) return NULL;

    Matrix* r = (Matrix*)malloc(sizeof(Matrix));
    r->cols = a->cols;
    r->rows = a->rows;
    float** v = (float**)malloc(sizeof(float*)*r->rows);

    for(int i = 0; i < r->rows; i++) {
        v[i] = (float*)malloc(sizeof(float)*r->cols);
        for(int k = 0; k < r->cols; k++) {
            v[i][k] = a->v[i][k] + b->v[i][k];
        }
    }
    r->v = v;
    return r;
}

Matrix* product(Matrix* a, Matrix* b) {

    if(a->cols != b->rows) return NULL;

    Matrix* r = (Matrix*)malloc(sizeof(Matrix));
    r->cols = b->cols;
    r->rows = a->rows;
    float** v = (float**)malloc(sizeof(float*)*r->rows);
    for(int i = 0; i < r->rows; i++) {
        v[i] = (float*)malloc(sizeof(float)*r->cols);
        for(int k = 0; k < r->cols; k++) {
            v[i][k] = 0;
            for(int j = 0; j < a->cols; j++) {
                v[i][k] += a->v[i][j] * b->v[j][k];
            }
        }
    }
    r->v = v;
    return r;
}

Matrix* transpose(Matrix* m) {

    Matrix* r = (Matrix*)malloc(sizeof(Matrix));
    r->cols = m->rows;
    r->rows = m->cols;
    float** v = (float**)malloc(sizeof(float*)*r->rows);
    for(int i = 0; i < r->rows; i++) {
        v[i] = (float*)malloc(sizeof(float)*r->cols);
        for(int k = 0; k < r->cols; k++) {
            v[i][k] = m->v[k][i];
        }
    }
    r->v = v;
    return r;
}

Matrix* layer_out(Layer* l, Matrix* input) {
    Matrix* p = product(l->weights, input);
    Matrix* r = sum(p, l->biases);
    free_matrix(p);
    return r;
}

Matrix* forward_pass(Network* nn, Matrix* input) {
    if(input == NULL) return NULL;
    if(nn->n_layers == 0) return NULL;
    if(input->cols != 1) return NULL; //needs to be a vector
    //if(input->rows != nn->layers[0].biases->rows) return NULL;

    Matrix* v = layer_out(&nn->layers[0], input);
    for(int i = 1; i < nn->n_layers; i++) {
        v = layer_out(&nn->layers[i], v);
    }
    return v;
}



// ----- Init functions -----

Layer* init_layer(int in_dim, int n_nodes) {

    Matrix* m = (Matrix*)malloc(sizeof(Matrix));
    m->rows = n_nodes;
    m->cols = in_dim;
    float** v = (float**)malloc(sizeof(float*)*m->rows);
    for(int i = 0; i < m->rows; i++) {
        v[i] = (float*)malloc(sizeof(float)*m->cols);
        for(int k = 0; k < m->cols; k++) {
            v[i][k] = _rand();
        }
    }
    m->v = v;

    Matrix* b = (Matrix*)malloc(sizeof(Matrix));
    b->rows = n_nodes;
    b->cols = 1;
    float** w = (float**)malloc(sizeof(float*)*b->rows);
    for(int i = 0; i < b->rows; i++) {
        w[i] = (float*)malloc(sizeof(float));
        w[i][0] = _rand();
    }
    b->v = w;

    Layer* l = (Layer*)malloc(sizeof(Layer));
    l->biases = b;
    l->weights = m;

    return l;
}

// ----- Free functions -----

void free_matrix(Matrix* m) {
    for(int i = 0; i < m->rows; i++) {
        free(m->v[i]);
    }
    free(m->v);
    free(m);
}

// ----- Print functions -----

void print_matrix(Matrix* m) {
    for(int i = 0; i < m->rows; i++) {
        printf("[");
        for(int k = 0; k < m->cols; k++) {
            printf(" %f", m->v[i][k]);
        }
        printf(" ]\n");
    }
    printf("\n");
}

void print_layer(Layer* l) {
    printf("biases: \n");
    print_matrix(l->biases);
    printf("weights: \n");
    print_matrix(l->weights);
}
