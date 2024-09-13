library IEEE;

use IEEE.std_logic_1164.all;
use IEEE.std_logic_textio.all;          -- For logic types I/O

library std;

use std.env.all;                        -- For hierarchical/external signals
use std.textio.all;                     -- For basic I/O

entity tb_TPU_MV_Element is
  generic(gCLK_HPER : time := 10 ns);  -- Generic for half of the clock cycle period
end tb_TPU_MV_Element;

architecture mixed of tb_TPU_MV_Element is
  component TPU_MV_Element is port(
      iCLK : in std_logic;
      iX : in integer;
      iW : in integer;
      iLdW : in integer;
      iY : in integer;
      oY : out integer;
      oX : out integer
      );
  end component;
-- Create signals for all of the inputs and outputs of the file that you are testing
-- := '0' or := (others => '0') just make all the signals start at an initial value of zero
  signal CLK, reset : std_logic := '0';
  signal s_iX : integer := 0;
  signal s_iW : integer := 0;
  signal s_iLdW : integer := 0;
  signal s_iY : integer := 0;
  signal s_oY : integer;
  signal s_oX : integer;
begin
  DUT0 : TPU_MV_Element
    port map(
      iCLK => CLK,
      iX => s_iX,
      iW => s_iW,
      iLdW => s_iLdW,
      iY => s_iY,
      oY => s_oY,
      oX => s_oX
    );
    
  P_CLK : process
  begin
    CLK <= '1';                         -- clock starts at 1
    wait for gCLK_HPER;                 -- after half a cycle
    CLK <= '0';                         -- clock becomes a 0 (negative edge)
    wait for gCLK_HPER;  -- after half a cycle, process begins evaluation again
  end process;
  
  P_RST : process
  begin
    reset <= '0';
    wait for gCLK_HPER/2;
    reset <= '1';
    wait for gCLK_HPER*2;
    reset <= '0';
    wait;
  end process;
  P_TEST_CASES : process
  begin
    
    wait for gCLK_HPER/2;  -- for waveform clarity, I prefer not to change inputs on clk edges
    -- Test case 1:
    -- Initialize weight value to 10.
    s_iX <= 0;  -- Not strictly necessary, but this makes the testcases easier to read
    s_iW <= 10;
    s_iLdW <= 1;
    s_iY <= 0;  -- Not strictly necessary, but this makes the testcases easier to read
    wait for gCLK_HPER*2;
    -- Expect: s_W internal signal to be 10 after positive edge of clock
    -- Test case 2:
    -- Perform average example of an input activation of 3 and a partial sum of 25. The weight is still 10. 
    s_iX <= 3;
    s_iW <= 0;
    s_iLdW <= 0;
    s_iY <= 25;
    wait for gCLK_HPER*2;
    wait for gCLK_HPER*2;
    -- Expect: o_Y output signal to be 55 = 3*10+25 and o_X output signal to 
    -- be 3 after two positive edge of clock.
    assert s_oY = 55 report "Test case 2 failed" severity error;
    -- Test case 3:
    -- Perform one MAC operation with minimum-case values
    s_iX <= 0;
    s_iW <= 10;
    s_iLdW <= 0;
    s_iY <= 0;
    wait for gCLK_HPER*2;
    wait for gCLK_HPER*2;
    -- Expect: o_Y output signal to be 0 = 10*0+0 and o_X output signal to be 0 
    -- after two positive edge of clock.
    assert s_oY = 0 report "Test case 3 failed" severity error;
    -- Test case 4:
    -- Change the weight and perform a MAC operation
    s_iX <= 2;
    s_iW <= 5;
    s_iLdW <= 1;
    s_iY <= 10;
    wait for gCLK_HPER*2;
    wait for gCLK_HPER*2;
    wait for gCLK_HPER*2;
    -- Expect: o_Y output signal to be 20 = 5*2+10 and o_X output signal 
    -- to be 2 after three positive edge of clock.
    assert s_oY = 20 report "Test case 4 failed" severity error;
    wait;
  end process;
  
end mixed;
