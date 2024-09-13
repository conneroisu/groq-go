library IEEE;
use IEEE.std_logic_1164.all;

entity mux2t1s is
  port (
    i_S  : in  std_logic;               -- selector
    i_D0 : in  std_logic;               -- data inputs
    i_D1 : in  std_logic;               -- data inputs
    o_O  : out std_logic                -- output
    );
end mux2t1s;

architecture structure of mux2t1s is

  component andg2 is
    port (
      i_A : in  std_logic;              -- input A to AND gate
      i_B : in  std_logic;              -- input B to AND gate
      o_F : out std_logic               -- output of AND gate
      );

  end component;

  component org2 is
    port (
      i_A : in  std_logic;              -- input A to OR gate
      i_B : in  std_logic;              -- input B to OR gate
      o_F : out std_logic               -- output of OR gate
      );

  end component;

  component invg is
    port (
      i_A : in  std_logic;              -- input to NOT gate
      o_F : out std_logic               -- output of NOT gate
      );

  end component;

  -- Signal to hold invert of the selector bit
  signal s_inv_S1   : std_logic;
  -- Signals to hold output valeus from 'AND' gates (needed to wire component to component?)
  signal s_oX, s_oY : std_logic;

begin
  ---------------------------------------------------------------------------
  -- Level 0: signals go through NOT stage
  ---------------------------------------------------------------------------
  invg1 : invg
    port map(
      i_A => i_S,                       -- input to NOT gate
      o_F => s_inv_S1                   -- output of NOT gate
      );
  ---------------------------------------------------------------------------
  -- Level 1: signals go through AND stage
  ---------------------------------------------------------------------------

  and1 : andg2
    port map(
      i_A => i_D0,                      -- input to AND gate
      i_B => s_inv_S1,                  -- input to AND gate
      o_F => s_oX                       -- output of AND gate
      );

  and2 : andg2
    port map(
      i_A => i_D1,                      -- input to AND gate
      i_B => i_S,                       -- input to AND gate
      o_F => s_oY                       -- output of AND gate
      );
  ---------------------------------------------------------------------------
  -- Level 1: signals go through OR stage (and then output)
  ---------------------------------------------------------------------------

  org1 : org2
    port map(
      i_A => s_oX,                      -- input to OR gate
      i_B => s_oY,                      -- input to OR gate
      o_F => o_O                        -- output of OR gate
      );
end structure;
