<script setup>
const socket = new WebSocket('ws://localhost:8080/socket')
socket.onmessage = function (event) {
  const data = JSON.parse(event.data)
  console.log(data)
}
</script>

<template>
  <h1>Test</h1>
</template>
