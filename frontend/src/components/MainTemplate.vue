<template>
  <div class="w3-container w3-khaki">
    <div class="main">
      <h1 class="w3-allerta">λόγος — Wiki for Literate Programming</h1>
      <hr>
      <p></p>
      <p></p>
      <footer class="w3-container w3-white" width="100%"></footer>
    </div>

    <div class="sidenav">
      <a href="/login">Login</a>
      <a href="/logout">Logout</a>
      <a href="/submit">New Page</a>
      <a href="/listsys">System List</a>
    </div>
  </div>
</template>

<script>
export default {
  name: 'MainTemplate',
  data() {
    return {
      siteTitle: 'λόγος — Wiki for Literate Programming',
      message: 'Your Message',
      port: '8080',
      addr: 'localhost'
    }
  }
}
</script>

<style>
@import url('https://www.w3schools.com/w3css/4/w3.css');
@import url('https://fonts.googleapis.com/css?family=Allerta+Stencil');

/* Add any additional CSS here */
</style>
