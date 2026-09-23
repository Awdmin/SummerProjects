#include <stdint.h>
#include <stdlib.h>
#include <unistd.h>
#include <stdio.h>
#include <string.h>
#include <sys/socket.h>
#include <netinet/in.h>
#include <linux/ip.h>
#include <linux/udp.h>


typedef struct _p_hdr {
    uint16_t src_port;
    uint16_t dst_port;
    uint16_t length;
    uint16_t checksum;
} p_hdr;

int main() {

    int sock = socket(AF_INET, SOCK_RAW, IPPROTO_RAW);

    return 0;
}
