<script>
  import Header from './lib/components/Header.svelte'
  import Calendar from './lib/components/Calendar.svelte'
  import Weather from './lib/components/Weather.svelte'
  import Tasks from './lib/components/Tasks.svelte'
</script>

<main>
  <Header />
  <div class="content">
    <div class="left-column">
      <Calendar />
    </div>
    <div class="right-column">
      <div class="right-top">
        <Weather />
      </div>
      <div class="right-bottom">
        <Tasks />
      </div>
    </div>
  </div>
</main>

<style>
  :global(body) {
    margin: 0;
    padding: 0;
    overflow: hidden;
  }

  :global(#app) {
    width: 100vw;
    height: 100vh;
    display: flex;
    flex-direction: column;
  }

  main {
    width: 100%;
    height: 100%;
    display: flex;
    flex-direction: column;
    background: #ffffff;
  }

  .content {
    width: 100%;
    height: 90%;
    display: flex;
    gap: 16px;
    padding: 16px;
    box-sizing: border-box;
    background: #f0f4f8;
  }

  .left-column {
    width: 60%;
    display: flex;
    flex-direction: column;
  }

  .right-column {
    width: 40%;
    display: flex;
    flex-direction: column;
    gap: 16px;
  }

  .right-top {
    width: 100%;
    height: 50%;
    display: flex;
  }

  .right-bottom {
    width: 100%;
    height: 50%;
    display: flex;
  }
</style>
