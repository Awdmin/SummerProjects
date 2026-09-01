#include <iso646.h>
#include <stdlib.h>
#include <stdio.h>
#include <math.h>

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

double _rand() {
    return rand() / (double)(RAND_MAX/2) - 1;
}

void a_func(Matrix* m, char id) {
    if(m->cols != 1) return;
    double s = 0;

    switch (id) {
        case 0: // sigmoid
            for(int i = 0; i < m->rows; i++) {
                m->v[i][0] = (1/(1 + pow(EULER_NUMBER, -m->v[i][0])));
            }
           break;
        case 1: // ReLU
            for(int i = 0; i < m->rows; i++) {
                m->v[i][0] = m->v[i][0] > 0 ? m->v[i][0] : 0;
            }
            break;
        case 2: // tanh
            for(int i = 0; i < m->rows; i++) {
                m->v[i][0] = tanh(m->v[i][0]);
            }
            break;
        case 3: // softmax
            for(int i = 0; i < m->rows; i++) {
                m->v[i][0] = pow(EULER_NUMBER, m->v[i][0]);
                s += m->v[i][0];
            }
            for(int i = 0; i < m->rows; i++) {
                m->v[i][0] /= s;
            }
            break;
        default:
            break;
    }
}

Matrix* sum(Matrix* a, Matrix* b) {

    if(a->cols != b->cols || a->rows != b->rows) return NULL;

    Matrix* r = (Matrix*)malloc(sizeof(Matrix));
    r->cols = a->cols;
    r->rows = a->rows;
    double** v = (double**)malloc(sizeof(double*)*r->rows);

    for(int i = 0; i < r->rows; i++) {
        v[i] = (double*)malloc(sizeof(double)*r->cols);
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
    double** v = (double**)malloc(sizeof(double*)*r->rows);
    for(int i = 0; i < r->rows; i++) {
        v[i] = (double*)malloc(sizeof(double)*r->cols);
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

void scalar(Matrix* m, double a) {
    for(int i = 0; i < m->rows; i++) {
        for(int k = 0; k < m->cols; k++) {
            m->v[i][k] *= a;
        }
    }
}

Matrix* transpose(Matrix* m) {

    Matrix* r = (Matrix*)malloc(sizeof(Matrix));
    r->cols = m->rows;
    r->rows = m->cols;
    double** v = (double**)malloc(sizeof(double*)*r->rows);
    for(int i = 0; i < r->rows; i++) {
        v[i] = (double*)malloc(sizeof(double)*r->cols);
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
    a_func(r, l->a_func);
    return r;
}

Matrix** forward_pass(Network* nn, Matrix* input) {
    if(input == NULL) return NULL;
    if(nn->n_layers == 0) return NULL;
    if(input->cols != 1) return NULL; //needs to be a vector

    Matrix** r = (Matrix**)malloc(sizeof(Matrix)*nn->n_layers);
    r[0] = layer_out(&nn->layers[0], input);
    for(int i = 1; i < nn->n_layers; i++) {
        r[i] = layer_out(&nn->layers[i], r[i-1]);
    }
    return r;
}

void backward_pass(Network* nn, Matrix** outputs) {

    Matrix* d_

}



// ----- Init functions -----

Layer* init_layer(int in_dim, int n_nodes, char a_func) {

    Matrix* m = (Matrix*)malloc(sizeof(Matrix));
    m->rows = n_nodes;
    m->cols = in_dim;
    double** v = (double**)malloc(sizeof(double*)*m->rows);
    for(int i = 0; i < m->rows; i++) {
        v[i] = (double*)malloc(sizeof(double)*m->cols);
        for(int k = 0; k < m->cols; k++) {
            v[i][k] = _rand();
        }
    }
    m->v = v;

    Matrix* b = (Matrix*)malloc(sizeof(Matrix));
    b->rows = n_nodes;
    b->cols = 1;
    double** w = (double**)malloc(sizeof(double*)*b->rows);
    for(int i = 0; i < b->rows; i++) {
        w[i] = (double*)malloc(sizeof(double));
        w[i][0] = _rand();
    }
    b->v = w;

    Layer* l = (Layer*)malloc(sizeof(Layer));
    l->biases = b;
    l->weights = m;
    l->a_func = a_func;

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

void print_network(Network* nn) {
    int t_params = 0;
    printf("-------- Neural Network info: (%d layers) --------\n", nn->n_layers);
    for(int i = 0; i < nn->n_layers; i++) {
        printf("L%d \t<> input dim: %d \t<> num nodes: %d\n", i, nn->layers[i].weights->cols, nn->layers[i].weights->rows);
        t_params += nn->layers[i].weights->cols * nn->layers[i].weights->rows + nn->layers[i].biases->rows;
    }
    printf("-------- Number of trainable params: %d -------\n", t_params);
}
